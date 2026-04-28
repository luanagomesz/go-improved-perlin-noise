package perlin

import (
	"math"
	"math/rand"
)

type (
	Perlin struct {
		seed        int64
		permutation []int
	}
	Vector2D struct {
		x, y float64
	}
	Vector3D struct {
		x, y, z float64
	}
)

func (v Vector2D) Dot(other Vector2D) float64 {
	return v.x*other.x + v.y*other.y
}

func (v Vector3D) Dot(other Vector3D) float64 {
	return v.x*other.x + v.y*other.y + v.z*other.z
}

func (p *Perlin) shufflePermutation(values []int) []int {
	r := rand.New(rand.NewSource(p.seed))
	for i := len(values) - 1; i > 0; i-- {
		randomIndex := r.Intn(i)
		values[i], values[randomIndex] = values[randomIndex], values[i]
	}
	return values
}

func (p *Perlin) makePermutation() []int {
	for i := 0; i < 256; i++ {
		p.permutation = append(p.permutation, i)
	}
	p.permutation = p.shufflePermutation(p.permutation)
	for i := 0; i < 256; i++ {
		p.permutation = append(p.permutation, p.permutation[i])
	}
	return p.permutation
}

func getConstantVector(v int) Vector2D {
	h := v % 4
	switch h {
	case 0:
		return Vector2D{1, 1}
	case 1:
		return Vector2D{-1, 1}
	case 2:
		return Vector2D{-1, -1}
	default:
		return Vector2D{1, -1}
	}
}

func getConstantVector1D(v int) float64 {
	if v%2 == 0 {
		return 1.0
	}
	return -1.0
}

func getConstantVector3D(v int) Vector3D {
	h := v % 12
	switch h {
	case 0:
		return Vector3D{1, 1, 0}
	case 1:
		return Vector3D{-1, 1, 0}
	case 2:
		return Vector3D{1, -1, 0}
	case 3:
		return Vector3D{-1, -1, 0}
	case 4:
		return Vector3D{1, 0, 1}
	case 5:
		return Vector3D{-1, 0, 1}
	case 6:
		return Vector3D{1, 0, -1}
	case 7:
		return Vector3D{-1, 0, -1}
	case 8:
		return Vector3D{0, 1, 1}
	case 9:
		return Vector3D{0, -1, 1}
	case 10:
		return Vector3D{0, 1, -1}
	default:
		return Vector3D{0, -1, -1}
	}
}

// fade applies the smoothstep curve 6t⁵ - 15t⁴ + 10t³ to eliminate first and second derivative
// discontinuities at integer boundaries, producing smoother noise output.
func fade(t float64) float64 {
	return ((6*t-15)*t + 10) * t * t * t
}

// lerp performs linear interpolation between a and b by factor t.
func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

// PerlinNoise2D returns the improved Perlin noise value for the 2D point (x, y).
// The result is in the approximate range [-1, 1].
func (p *Perlin) PerlinNoise2D(x, y float64) float64 {
	px := int(math.Floor(x)) % 256
	py := int(math.Floor(y)) % 256

	xf := x - math.Floor(x)
	yf := y - math.Floor(y)

	topRight := Vector2D{xf - 1.0, yf - 1.0}
	topLeft := Vector2D{xf, yf - 1.0}
	bottomRight := Vector2D{xf - 1.0, yf}
	bottomLeft := Vector2D{xf, yf}

	valueTopRight := p.permutation[p.permutation[px+1]+py+1]
	valueTopLeft := p.permutation[p.permutation[px]+py+1]
	valueBottomRight := p.permutation[p.permutation[px+1]+py]
	valueBottomLeft := p.permutation[p.permutation[px]+py]

	dotTopRight := topRight.Dot(getConstantVector(valueTopRight))
	dotTopLeft := topLeft.Dot(getConstantVector(valueTopLeft))
	dotBottomRight := bottomRight.Dot(getConstantVector(valueBottomRight))
	dotBottomLeft := bottomLeft.Dot(getConstantVector(valueBottomLeft))

	u := fade(xf)
	v := fade(yf)

	return lerp(
		lerp(dotBottomLeft, dotTopLeft, v),
		lerp(dotBottomRight, dotTopRight, v),
		u,
	)
}

// Noise2D returns 2D fractional Brownian motion (fBm) noise by summing multiple octaves of PerlinNoise2D.
// Each octave doubles the frequency and halves the amplitude, adding progressively finer detail.
func (p *Perlin) Noise2D(x, y float64, numOctaves int, frequency float64) float64 {
	result := 0.0
	amplitude := 1.0

	for range numOctaves {
		result += amplitude * p.PerlinNoise2D(x*frequency, y*frequency)
		amplitude *= 0.5
		frequency *= 2.0
	}

	return result
}

// PerlinNoise1D returns the improved Perlin noise value for the 1D point x.
// The result is in the approximate range [-1, 1].
func (p *Perlin) PerlinNoise1D(x float64) float64 {
	px := int(math.Floor(x)) & 255
	xf := x - math.Floor(x)

	left := xf * getConstantVector1D(p.permutation[px])
	right := (xf - 1.0) * getConstantVector1D(p.permutation[px+1])

	return lerp(left, right, fade(xf))
}

// Noise1D returns 1D fractional Brownian motion (fBm) noise by summing multiple octaves of PerlinNoise1D.
// Useful for animating values over time or generating simple terrain profiles.
func (p *Perlin) Noise1D(x float64, numOctaves int, frequency float64) float64 {
	result := 0.0
	amplitude := 1.0

	for range numOctaves {
		result += amplitude * p.PerlinNoise1D(x*frequency)
		amplitude *= 0.5
		frequency *= 2.0
	}

	return result
}

// PerlinNoise3D returns the improved Perlin noise value for the 3D point (x, y, z).
// The result is in the approximate range [-1, 1].
func (p *Perlin) PerlinNoise3D(x, y, z float64) float64 {
	px := int(math.Floor(x)) & 255
	py := int(math.Floor(y)) & 255
	pz := int(math.Floor(z)) & 255

	xf := x - math.Floor(x)
	yf := y - math.Floor(y)
	zf := z - math.Floor(z)

	aaa := p.permutation[p.permutation[p.permutation[px]+py]+pz]
	aba := p.permutation[p.permutation[p.permutation[px]+py+1]+pz]
	aab := p.permutation[p.permutation[p.permutation[px]+py]+pz+1]
	abb := p.permutation[p.permutation[p.permutation[px]+py+1]+pz+1]
	baa := p.permutation[p.permutation[p.permutation[px+1]+py]+pz]
	bba := p.permutation[p.permutation[p.permutation[px+1]+py+1]+pz]
	bab := p.permutation[p.permutation[p.permutation[px+1]+py]+pz+1]
	bbb := p.permutation[p.permutation[p.permutation[px+1]+py+1]+pz+1]

	dot := func(hash int, dx, dy, dz float64) float64 {
		return getConstantVector3D(hash).Dot(Vector3D{dx, dy, dz})
	}

	x1 := lerp(dot(aaa, xf, yf, zf), dot(baa, xf-1, yf, zf), fade(xf))
	x2 := lerp(dot(aba, xf, yf-1, zf), dot(bba, xf-1, yf-1, zf), fade(xf))
	y1 := lerp(x1, x2, fade(yf))

	x3 := lerp(dot(aab, xf, yf, zf-1), dot(bab, xf-1, yf, zf-1), fade(xf))
	x4 := lerp(dot(abb, xf, yf-1, zf-1), dot(bbb, xf-1, yf-1, zf-1), fade(xf))
	y2 := lerp(x3, x4, fade(yf))

	return lerp(y1, y2, fade(zf))
}

// Noise3D returns 3D fractional Brownian motion (fBm) noise by summing multiple octaves of PerlinNoise3D.
// Suitable for volumetric textures, cave generation, or biome density fields.
func (p *Perlin) Noise3D(x, y, z float64, numOctaves int, frequency float64) float64 {
	result := 0.0
	amplitude := 1.0

	for range numOctaves {
		result += amplitude * p.PerlinNoise3D(x*frequency, y*frequency, z*frequency)
		amplitude *= 0.5
		frequency *= 2.0
	}

	return result
}

// NewPerlin creates and returns a new Perlin noise generator initialized with the given seed.
func NewPerlin(seed int64) *Perlin {
	p := &Perlin{seed: seed}
	p.makePermutation()
	return p
}
