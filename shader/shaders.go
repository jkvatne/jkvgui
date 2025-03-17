package shader

var FragmentQuadShader = `#version 400
in vec2 fragTexCoord;
out vec4 outputColor;

uniform sampler2D tex;
uniform vec4 textColor;

void main() {    
    vec4 sampled = vec4(1.0, 1.0, 1.0, texture(tex, fragTexCoord).r);
    outputColor = textColor * sampled;
}	
` + "\x00"

var VertexQuadShader = `#version 400
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

var RectVertShaderSource = `
	#version 330
	layout(location = 1) in vec2 inPos;
	uniform vec2 resolution;

	void main() {
		vec2 clipSpace = (inPos / resolution) * 2.0 - 1.0;
		gl_Position = vec4(clipSpace * vec2(1, -1), 0, 1);
	}
	` + "\x00"

var RectFragShaderSource = `
	#version 330
	in vec4 gl_FragCoord;
	out vec4 fragColor;

	uniform vec2 pos;
	uniform vec2 halfbox;
    uniform vec2 rw;  // Corner radius, border width
	uniform vec4 colors[2]; // Fillcolor, FrameColor
	uniform vec2 resolution;

	float sdRoundedBox( in vec2 p, in vec2 b, in float r ) {
		vec2 q = abs(p)-b+r;
		return min(max(q.x,q.y),0.0) + length(max(q,0.0)) - r;
	}

	void main() {
		fragColor = colors[1]; // FrameColor
		float bw = rw.y;  // Frame width + shadow width
        float rr = rw.x;  // Corner radius
        vec2 p = vec2(gl_FragCoord.x-pos.x, resolution.y-gl_FragCoord.y-pos.y);
		// halfbox includes shadow. hb1 subtracts shadow to get frame size.
        // Now d1 is distance from frame
		float d1 = sdRoundedBox(p, halfbox, rr);

		// hb2 is the inside of the frame.
		vec2 hb2 = vec2(halfbox.x-bw, halfbox.y-bw);
		float d2 = sdRoundedBox(p, hb2, rr-rw.y);
		if (d1>0.0) {
			vec4 col = vec4(0.0, 0.0, 0.0, 0.0);
			fragColor = mix(colors[1], col, clamp(d1, 0, 1));
		}
		if (d2<=0.5) { 
            // We are inside box. Mix with border to smooth border
			fragColor = mix(colors[1], colors[0], clamp(1-d2, 0, 1));
		}
	}
	` + "\x00"

var ShadowFragShaderSource = `
	#version 330
	in vec4 gl_FragCoord;
	out vec4 fragColor;

	uniform vec2 pos;
	uniform vec2 halfbox;
    uniform vec4 rws;  // Corner radius, border width, shaddow size, shadow alfa
	uniform vec2 resolution;

	float sdRoundedBox( in vec2 p, in vec2 b, in float r ) {
		vec2 q = abs(p)-b+r;
		return min(max(q.x,q.y),0.0) + length(max(q,0.0)) - r;
	}

	void main() {
		float sw = rws.z;  // Shadow width
        float rr = rws.x;  // Corner radius
        fragColor = vec4(0.3,0.3, 0.3, 0.3);
        vec2 p = vec2(gl_FragCoord.x-pos.x, resolution.y-gl_FragCoord.y-pos.y);
		// halfbox includes shadow. hb1 subtracts shadow to get frame size.
        vec2 hb1 = vec2(halfbox.x-sw, halfbox.y-sw);
        // Now d1 is distance from shadow center
		float d1 = sdRoundedBox(p, hb1, rr);
		if (d1>-sw) {
			// Outside frame
            float alfa = 0.3 * smoothstep(0,-sw,d1);
			fragColor = vec4(0.3, 0.3, 0.3, max(0.0, alfa));	
		}
	}
	` + "\x00"
