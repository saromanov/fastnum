# Ресёрч: свежие публикации и их отношение к fastnum

Сессия 2026-11-09. Поиск: arXiv API (`export.arxiv.org`), Go 1.27 release notes,
pkg.go.dev. Контекст: fill константами (zeros/ones) + проверка массивов на
константу; i7-9750H (6C/12T, L1d 32K, L2 256K, L3 12M), Go 1.27.3 amd64.

## Найденные работы

### Прямо релевантные

1. **arXiv 2311.04797 — «CloverLeaf on Intel Multi-Core CPUs: A Case Study
   in Write-Allocate Evasion» (2023).**
   Ключевая работа для нашего случая. При обычной записи каждая кэш-линия
   сначала читается из DRAM (read-for-ownership, write-allocate), т.е. чистая
   запись стоит примерно 2× пропускной способности DRAM. Обход:
   non-temporal (streaming) stores, запись полными кэш-линиями.
   Подтверждено нашими числами: memclr (пишет полными линиями) — 21–26 GB/s;
   копирующий ones-fill (seed-чтение + RFO) — ~10 GB/s.

2. **arXiv 2409.08108 — «Microarchitectural comparison and in-core modeling
   of state-of-the-art CPUs: Grace, Sapphire Rapids…» (2024).**
   In-core модели пропускной способности памяти; подтверждает потолки,
   которые мы наблюдали эмпирически.

3. **arXiv 1910.05109 — «Base64 encoding and decoding at almost the speed
   of a memory copy» (Lemire et al., 2019).**
   Методология «держаться на пределе memcpy». Наш сдвиговый compare
   (`bytes.Equal(b[8:], b[:len-8])`) — та же идея: один поток чтения вместо
   двух.

4. **arXiv 2610.17406 — «Fast Forgiving Base64 Decoding with SIMD
   Instructions» (окт 2026).** Свежейшая линия SIMD-обработки данных:
   ручные интринсики вместо надежд на автогенерацию.

5. **arXiv 2412.16370 — «Faster Positional-Population Counts for AVX2,
   AVX-512, and ASIMD» (2024).** Паттерн широкорегистровых редукций
   (accum + VPTEST) — применим к проверкам, если когда-нибудь понадобится.

6. **arXiv 2611.03917 — «Silent Degradation during Auto-Vectorization of
   bf16/fp16 on RISC-V Under LLVM» (ноя 2026).** Компиляторы молча теряют
   векторизацию → ручной SIMD обязателен; аргумент против «ручных
   разверток» на чистом Go.

### Поиск дал мало (что тоже информация)

Запросы по точным терминам `memset`, `memcpy optimization`, `non-temporal
stores` на arXiv практически пусты: такие работы публикуются в ACM/USENIX
(пример практического стандарта — реализация memset/memcpy в glibc 2.3x
с EVEX/NT-путями, Intel Optimization Manual §«Non-Temporal Stores»).
Semantic Scholar API отдает 403 (Forbidden) без ключа.

## Платформа: Go 1.27

- `GOEXPERIMENT=simd` → два пакета: переносимый `simd` и архитектурный
  `simd/archsimd` (начат в 1.26; в 1.27 пересмотрен amd64-API, добавлены
  arm64 NEON и wasm). Типы `Int64x8`/`Float64x8` (512-bit), `Broadcast*`,
  `Load*/Store*`, `Mask*`. API объявлен нестабильным.
- **NT-stores в API archsimd нет** (проверено по списку функций pkg.go.dev).
- Рантайм-примитивы (`memclr`, `memmove`, `memequal`) уже используют AVX2
  при `GOAMD64=v3` (на нашем CPU v3 дает: fill −18%, is-проверки +12..25%
  хуже — ERMS `rep movsb/cmpsb` на Coffee Lake быстрее AVX2-петель).

## Измерения, относящиеся к выводам (см. BENCHMARKS.md)

- Параллельный fill 128 MB (2/4/12 горутин): прирост в шуме (11.8 → 10.8 ms
  на 4 потоках) — память упёрта в полосу записи.
- Zeros 16M: 5.1 ms (`GOAMD64=v3`) = ~26 GB/s ≈ крыша записи этой машины.

## Вердикт: что реализовывать

1. **Стоит: NT-stores для ones-fill на буферах > L3** (единственное место с
   запасом >2×). Реализация: Go-asm (удобно генерировать avo), выбор
   AVX2/AVX-512 через `golang.org/x/sys/cpu`, чисто-Go fallback, порог
   включения (например > 4 MB: ниже NT store вредит — eviction без выгоды).
   Прогноз: ones 16M ~ 7 ms вместо 11–13.
2. **Не стоит сейчас:** `simd/archsimd` (не для скорости: рантайм уже AVX2,
   NT-операций нет, API нестабилен); параллелизм (измерено — нет эффекта на
   128 MB); дальнейшие pure-Go микродвижения (все упираются в память).
