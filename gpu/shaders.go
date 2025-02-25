package gpu

var (
	RectFragShaderSource = `
	#version 330

	uniform vec2 pos;
	uniform vec2 halfbox;
    uniform vec2 rw;
	uniform vec4 colors[2];
	uniform vec2 resolution;

	in vec4 gl_FragCoord;

	out vec4 fragColor;

	float sdRoundedBox( in vec2 p, in vec2 b, in float r ) {
		vec2 q = abs(p)-b+r;
		return min(max(q.x,q.y),0.0) + length(max(q,0.0)) - r;
	}

	void main() {
		fragColor = colors[1];
        vec2 p = vec2(gl_FragCoord.x-pos.x, resolution.y-gl_FragCoord.y-pos.y);
		float d1 = sdRoundedBox(p, halfbox, rw.x);
		vec2 hb2 = vec2(halfbox.x-rw.y*2, halfbox.y-rw.y*2);
		float d2 = sdRoundedBox(p, hb2, rw.x-rw.y);
		if (d1>0.0) {
			discard;
		}
		if (d2<=0) {
			fragColor = colors[0];
		}
	}
	` + "\x00"

	RectVertShaderSource = `
	#version 330
	layout(location = 1) in vec2 inPos;
	uniform vec2 resolution;
	void main() {
		vec2 zeroToOne = inPos / resolution;
		vec2 zeroToTwo = zeroToOne * 2.0;
		vec2 clipSpace = zeroToTwo - 1.0;
		gl_Position = vec4(clipSpace * vec2(1, -1), 0, 1);
	}
	` + "\x00"
)
