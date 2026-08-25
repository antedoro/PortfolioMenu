# PortfolioMenu

PortfolioMenu e' una piccola app Go per macOS che vive nella menu bar e mostra
una dashboard locale del portafoglio su `http://localhost:8080`.

## Funzioni

- Aggiornamento periodico dei prezzi da provider esterni.
- Dashboard web con riepilogo, allocazione, storico e fiscalita' italiana.
- Gestione asset da interfaccia: aggiunta, modifica, duplicazione ed eliminazione.
- Esportazione CSV del portafoglio.
- Tema chiaro/scuro salvato nella configurazione.

## Avvio in sviluppo

```sh
go run ./cmd/portfoliomenu
```

La configurazione principale e' in `configs/portfoliomenu.toml`. Il portafoglio
viene salvato nel file indicato da `data_file`, di default
`configs/portfolio.toml`.

## Test

```sh
GOCACHE="$PWD/.gocache" go test ./...
```

La cache locale evita di scrivere nella cache globale di Go quando l'ambiente e'
sandboxato.
