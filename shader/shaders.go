package shader

var FragmentFontShader = `#version 400
in vec2 fragTexCoord;
out vec4 outputColor;

uniform sampler2D tex;
uniform vec4 textColor;

void main() {    
    vec4 sampled = vec4(1.0, 1.0, 1.0, texture(tex, fragTexCoord).r);
    outputColor = textColor * sampled;
}	
` + "\x00"

var VertexFontShader = `#version 400
in vec2 vert;
in vec2 vertTexCoord;
out vec2 fragTexCoord;

uniform vec2 resolution;

void main() {
   vec2 clipSpace = (vert / resolution) * 2.0 - 1.0;
   gl_Position = vec4(clipSpace * vec2(1, -1), 0, 1);
   fragTexCoord = vertTexCoord;
}
` + "\x00"

var RectFragShaderSource = `
	#version 330
	in vec4 gl_FragCoord;
	out vec4 fragColor;

	uniform vec2 pos;
	uniform vec2 halfbox;
    uniform vec2 rw;
	uniform vec4 colors[2];
	uniform vec2 resolution;

	float sdRoundedBox( in vec2 p, in vec2 b, in float r ) {
		vec2 q = abs(p)-b+r;
		return min(max(q.x,q.y),0.0) + length(max(q,0.0)) - r;
	}

	void main() {
		fragColor = colors[1];
		float bw = 11.0;  // Border width
		float sw = 10.0;  // Shadow width
        float rr = 10.0; // Corner radius
        vec2 p = vec2(gl_FragCoord.x-pos.x, resolution.y-gl_FragCoord.y-pos.y);
        vec2 hb1 = vec2(halfbox.x-sw, halfbox.y-sw);
		float d1 = sdRoundedBox(p, hb1, rr);

		vec2 hb2 = vec2(halfbox.x-bw, halfbox.y-bw);
		float d2 = sdRoundedBox(p, hb2, rr);
		if (d1>0.0) {
			fragColor = vec4(0.6, 0.6, 0.6, max(0.0,1.0-min(1.0,1.4*d1/sw)));
		}
		if (d2<=0) {
			fragColor = colors[0];
		}
	}
	` + "\x00"

var RectVertShaderSource = `
	#version 330
	layout(location = 1) in vec2 inPos;
	uniform vec2 resolution;

	void main() {
		vec2 clipSpace = (inPos / resolution) * 2.0 - 1.0;
		gl_Position = vec4(clipSpace * vec2(1, -1), 0, 1);
	}
	` + "\x00"
