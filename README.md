# go-improved-perlin-noise

A Go implementation of Ken Perlin's improved noise algorithm (2002), supporting 1D, 2D, and 3D noise with or without fractional Brownian motion (fBm)

Each dimension comes in two variants:

- **`PerlinNoise`** — classic single-sample Perlin noise. Returns the raw gradient noise value for a given point. Use this when you want direct control over how the noise is sampled.
- **`Noise`** — fractional Brownian motion (fBm) built on top of Perlin noise. Layers multiple octaves of noise at increasing frequencies and decreasing amplitudes, producing richer and more natural-looking results. Use this for easier better results

## Installation

```bash
go get github.com/luanagomesz/go-improved-perlin-noise
```

## Usage

```go
import perlin "github.com/luanagomesz/go-improved-perlin-noise"

p := perlin.NewPerlin(42) // provide a seed
```

### 1D Noise

```go
// Classic Perlin noise
value := p.PerlinNoise1D(x)

// fBm — multiple octaves layered together
value := p.Noise1D(x, octaves, frequency)
```

### 2D Noise

```go
// Classic Perlin noise
value := p.PerlinNoise2D(x, y)

// fBm — multiple octaves layered together
value := p.Noise2D(x, y, octaves, frequency)
```

### 3D Noise

```go
// Classic Perlin noise
value := p.PerlinNoise3D(x, y, z)

// fBm — multiple octaves layered together
value := p.Noise3D(x, y, z, octaves, frequency)
```

## Parameters

| Parameter   | Description                                                              |
|-------------|--------------------------------------------------------------------------|
| `seed`      | Integer seed for reproducible results                                    |
| `x/y/z`     | Floating-point coordinates to sample                                     |
| `octaves`   | Number of noise layers — more octaves add finer detail (fBm only)        |
| `frequency` | Initial frequency — higher values produce denser noise (fBm only)        |

## Output Range

functions return values in the approximate range `[-1, 1]`.

## References

- [Perlin Noise: A Procedural Generation Algorithm — Raouf Touti](https://rtouti.github.io/graphics/perlin-noise-algorithm)
- [Improved Noise reference implementation — Ken Perlin](https://mrl.cs.nyu.edu/~perlin/noise/)

## License

[MIT](LICENSE)
