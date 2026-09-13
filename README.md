# molinets

[![Unit Tests](https://github.com/gllm-dev/molinets/actions/workflows/test.yml/badge.svg)](https://github.com/gllm-dev/molinets/actions/workflows/test.yml)
[![License: GPL v3](https://img.shields.io/badge/license-GPLv3-blue.svg)](LICENSE)

An MCP server that translates coffee grinder settings between 34 popular manual and electric grinders.

*Molinets* is Catalan for "grinders".

A recipe says "Comandante, 26 clicks" and you own a Timemore? Ask your assistant. It calls molinets and answers in your grinder's own notation:

```text
> I'm at 26 clicks on a Comandante. What's that on my C3S Pro, a Kingrinder K6 and a Fellow Ode?

  ≈ 710 µm, pour-over to French press territory
  Timemore C3S Pro   19
  Kingrinder K6      1.24     (1 rotation, 24 clicks)
  Fellow Ode Gen 2   6
```

It can also give you a starting range for a brew method on your grinder, or the setting for a target particle size.

## Quick start

molinets runs anywhere Docker does. The image is multi-arch (amd64 and arm64).

**As a local server (stdio).** Add this to your MCP client's configuration. Most desktop clients use this `mcpServers` format:

```json
{
  "mcpServers": {
    "molinets": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "ghcr.io/gllm-dev/molinets:latest", "stdio"]
    }
  }
}
```

**As a remote server (Streamable HTTP).**

```sh
docker run -d --name molinets -p 8080:8080 ghcr.io/gllm-dev/molinets:latest
```

Point your client at `http://localhost:8080/mcp`, or at `https://<your-host>/mcp` behind a reverse proxy. The transport is stateless, so it scales horizontally without sticky sessions. `GET /healthz` answers `{"status":"ok"}`.

There is no authentication. Every tool is read-only and the data is public.

## Tools

| Tool | Arguments | What you get |
| --- | --- | --- |
| `convert_grind_setting` | `from_grinder`, `setting`, `to_grinders` *(optional, all by default)* | Approximate particle size, the brew methods that suit it, and the equivalent setting on each target grinder |
| `grind_setting_for_method` | `grinder`, `method` | A from / middle / to setting range for that method on that grinder |
| `microns_to_setting` | `grinder`, `microns` | The setting that produces that particle size |
| `list_grinders` | `type` *(optional: `manual` or `electric`)* | Every supported grinder with its aliases, notation, range and microns per step |

**Grinder names.** You don't need exact names: `comandante`, `c3s pro`, `k6`, `ode`, `niche` and `Mahlkoenig X54` all resolve. An ambiguous name such as `hario` returns the candidates instead of guessing.

**Brew methods.** You can name a method (`espresso`, `v60`, `french press`) or a brewer (`chemex`, `origami`, `switch`, `clever`, `moccamaster`).

**Errors.** Failed calls start with a stable code, so clients and models can react to them:

- `GRINDER_NOT_FOUND`
- `AMBIGUOUS_GRINDER`
- `INVALID_SETTING`
- `SETTING_OUT_OF_RANGE`
- `METHOD_NOT_FOUND`
- `METHOD_UNREACHABLE`
- `INVALID_ARGUMENT`

### Setting notation

Settings are written the way each grinder's dial reads:

| Notation | Example | Used by |
| --- | --- | --- |
| clicks or steps | `26` | Comandante, Timemore C2 and C3, Kingrinder P series, Baratza, DF64, Niche Zero (`53` means the second rotation) |
| rotation.number.click | `2.8.4` | 1Zpresso J, K, X and Q series, Timemore C3 ESP Pro |
| rotation.clicks | `1.22` (82 clicks) | Kingrinder K4 and K6 |
| number.click | `4.5` | 1Zpresso ZP6, Timemore C5 Pro and S3 |
| number.thirds | `4.2` (4⅔) | Fellow Ode Gen 2, OXO |
| decimal dial | `5.25`, `12.5` | Fellow Opus, Mahlkönig X54, Timemore Sculptor, Varia VS3 |

Settings that can't exist on a grinder are rejected, with the valid range in the message. For example, `1.9.0` is rejected on a J-Max because it has 9 numbers per rotation, and 60 clicks are rejected on a Comandante.

## Supported grinders

**Manual:**

- Comandante C40 MK4 (with and without Red Clix)
- 1Zpresso J-Max, J-Ultra, K-Ultra, K-Max (also K-Plus and K-Pro), ZP6 Special, X-Pro / X-Ultra, Q2 / Q Air
- Timemore C2, C3 / C3S / C3S Pro, C3 ESP Pro, C5 Pro, S3
- Kingrinder K4 / K6, P1 / P2
- Hario Skerton Pro, Mini Mill Plus
- Porlex Mini II

**Electric:**

- Baratza Encore, Encore ESP, Virtuoso+, Sette 270
- Fellow Ode Gen 2, Opus
- Niche Zero
- Turin DF64 Gen 2
- Sage / Breville Smart Grinder Pro
- Wilfa Svart Aroma, Uniform
- Timemore Sculptor 064S
- Varia VS3 Gen 2
- OXO Conical Burr
- Mahlkönig X54

Missing yours? See [Contributing](#contributing).

## How it works

Every conversion goes through an approximate particle size:

```text
setting on grinder A  →  microns  →  setting on grinder B
```

Each grinder has a calibration from its smallest step to microns. The calibrations are fitted to the per-grinder brew tables of the [Honest Coffee Guide](https://honestcoffeeguide.com/coffee-grind-size-chart/), and the test suite checks all 34 grinders against 263 points from those tables.

Manufacturer "microns per click" figures are deliberately ignored. On conical burrs they usually measure how far the burr moves, not how much the grounds change, and they overstate it two to three times. `list_grinders` shows both figures side by side.

The results are starting points to dial in by taste, for two reasons:
- Two units of the same model can find zero one or two clicks apart. Count from burr contact on yours.
- A grind is a spread of particle sizes, not one number.

## Running from source

Requires Go 1.24 or later.

```sh
git clone https://github.com/gllm-dev/molinets.git
cd molinets
go build -o molinets cmd/main.go

./molinets http     # Streamable HTTP on :8080
./molinets stdio    # stdin/stdout, for local clients
```

| Setting | Default | |
| --- | --- | --- |
| `PORT` | `8080` | HTTP port |
| `--log-level` | `info` | `debug`, `info`, `warn` or `error`. In stdio mode, logs go to stderr. |

## Contributing

Issues and pull requests are welcome, especially new grinders and better data for existing ones.

To add a grinder:
1. Add it to [`internal/adapters/repositories/embedded/data/grinders.json`](internal/adapters/repositories/embedded/data/grinders.json) with its notation, calibration and range.
2. Add its reference settings to [`testdata/hcg_reference.json`](internal/adapters/repositories/embedded/testdata/hcg_reference.json).
3. Run the tests. They check the new grinder against its reference points, and check that its settings format cleanly and convert both ways with every other grinder.

```sh
go vet ./...
go test -race ./...
```

Dependency injection is generated with [Wire](https://github.com/google/wire). After changing a constructor, run `cd di && wire`. CI fails if `di/wire_gen.go` is out of date.

## License

[GNU General Public License v3.0](LICENSE).
