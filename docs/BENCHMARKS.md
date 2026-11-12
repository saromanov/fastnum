# Методика и результаты бенчмарков

2026-11-09. Машина: i7-9750H (6C/12T), Go 1.27.3 linux/amd64, по умолчанию
`GOAMD64=v1`. Прогоны: `-count=6 -benchtime=500ms -benchmem`, сравнение
через `benchstat` (медианы, p-value). Прогрев между прогонами не изолирован,
поэтому разброс до ±30% на мелких размерах — мелкие числа читать с осторожностью.

## Сравниваемые реализации

- **impl** — публичный API пакета (текущая/экспериментальная версия).
- **loop** — «честный» рукописный цикл: `slices.Grow` + fill/`clear`.
- **naive** — наивный `for { dst = append(dst, v) }`.

Размеры: 8, 64, 1K, 64K, 16M элементов (до 128 MB).

## Фазы эксперимента

1. `docs/benchmarks/old.txt` — исходная реализация (таблицы 64 KB, append-копирование).
2. `docs/benchmarks/new.txt` — первый вариант правок «как в совете»:
   `append(dst, make(T, n)...)` для zeros + чистое удвоение seed для ones +
   сдвиговый compare. Регрессии на мелких размерах: AppendOnes/64 +139%,
   /1K +65%, AppendZeros/8 +32%.
3. `docs/benchmarks/final.txt` —zeros = `Grow`+`clear`, ones = seed-таблица 8 KB +
   стриминг из горячего префикса, is = сдвиговый compare. Осталась регрессия
   AppendOnes/8 (+46%).
4. `docs/benchmarks/final2.txt` — + быстрый путь `append(dst, seed[:items]...)`
   при items ≤ 1024. Регрессия снята (в пределах шума, ~1 ns).
5. `docs/benchmarks/cur-v1.txt`, `cur-v3.txt` — `GOAMD64=v1` vs `v3` на текущем коде
   (размеры 64K/16M): fills −11..−19%, is-проверки +12..+26% хуже.
6. `BenchmarkParallelFillOnes` (par_fill_test.go): 128 MB, 2/4/12 горутин —
   прирост в шуме (11.8 → 10.8 ms), память упёрта в полосу записи.

## Итоговые числа (impl, старая → финальная, медианы)

| Бенчмарк | Было | Стало | Δ |
|---|---|---|---|
| AppendInt64Zeros/16M | 13.0 ms | 6.2 ms | **−52%** |
| AppendFloat64Zeros/16M | 15.6 ms | 6.5 ms | **−58%** |
| AppendInt64Ones/16M | 27.7 ms | 13.7 ms | **−50%** |
| AppendInt64Ones/64K | 24.3 µs | 16.8 µs | −31% |
| AppendFloat64Zeros/64K | 17.3 µs | 12.3 µs | −29% |
| AppendFloat64Zeros/64 | 16.8 ns | 12.7 ns | −24% |
| AppendInt64Zeros/64 | 14.5 ns | 10.7 ns | −26% |
| IsInt64Ones/8 | 14.3 ns | 8.6 ns | **−40%** |
| IsFloat64Zeros/8 | 15.8 ns | 9.7 ns | −38% |
| IsFloat64Zeros/64 | 36.4 ns | 25.0 ns | −31% |
| IsInt64Zeros/8 | 12.5 ns | 10.4 ns | −17% |
| IsInt64Ones/64K | 22.7 µs | 16.1 µs | −29% |
| IsInt64Ones/16M | 13.2 ms | 9.8 ms | −26% |
| IsFloat64Ones/16M | 12.0 ms | 9.3 ms | −23% |
| IsFloat64Ones/1K | 210 ns | 158 ns | −25% |

Geomean всей нагрузки: −6.3%. Аллокации: 0 везде (буфер с готовым cap).

Изменённые реализации после всех правок:

- **zeros:** оставлена исходная таблица 64 KB: `clear` не дал стабильного
  преимущества на основном размере 8K и был исключён из минимального PR.
- **ones:** ≤16 — обычный цикл для int64 и `append` для float64;
  17–1024 — один `append` из таблицы-семечка (8 KB, L1); больше — первый chunk
  из таблицы, далее стриминг `copy(s[f:], s[:1024])` из горячего префикса.
- **is:** ≤16 — простой цикл; больше — `bytes.Equal(b[8:], b[:len(b)-8])`
  (сдвиговое сравнение: все элементы равны) + проверка `a[0]`; float64
  побитово (`math.Float64bits`), семантика −0/NaN сохранена со старой версией.

## Верхние границы (roofline по нашим данным)

- Запись DRAM: ~26 GB/s (zeros 16M, v3: 4.58 ms / 128 MB).
- Запись с RFO (ones copy): ~10 GB/s → именно здесь запас для NT-stores
  (см. RESEARCH.md).
- Чтение DRAM (is-проверки): ~16 GB/s, близко к крыше чтения.
7. `docs/benchmarks/v0.1.0.txt` — прогон релиза v0.1.0 (`-benchtime=100ms
   -count=3 -benchmem`, GOAMD64=v1, Go 1.27.3); источник чисел в README.
