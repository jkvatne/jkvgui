package gpu

// FragQuadSource is a fragment shader that draws a rectangle with texture. Used by fonts and icons.
var FragQuadSource = `#version 330
in vec2 fragTexCoord;
out vec4 outputColor;

uniform sampler2D tex;
uniform vec4 textColor;

void main() {    
    vec4 sampled = vec4(1.0, 1.0, 1.0, texture(tex, fragTexCoord).r);
    outputColor = textColor * sampled;
}	
` + "\x00"

// VertQuadSource is a vertex shader that draws a rectangle with texture. Used by fonts and icons.
var VertQuadSource = `#version 330
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

// FragImgSource is the fragment shader used to draw images.
var FragImgSource = `#version 330
in vec2 fragTexCoord;
out vec4 outputColor;

uniform sampler2D tex;
uniform vec4 textColor;

void main() {    
    outputColor = texture(tex, fragTexCoord); 
}	
` + "\x00"

// VertRectSource is a vertex shader used for rounded rectangles and shaddows
var VertRectSource = `
	#version 330
	layout(location = 1) in vec2 inPos;
	uniform vec2 resolution;

	void main() {
		vec2 clipSpace = (inPos / resolution) * 2.0 - 1.0;
		gl_Position = vec4(clipSpace * vec2(1, -1), 0, 1);
	}
	` + "\x00"

// FragRectSource is the fragment shader used for rounded rectangles
var FragRectSource = `
	#version 330
	in vec4 gl_FragCoord;
	out vec4 fragColor;

	uniform vec2 pos;
	uniform vec2 halfbox;
    uniform vec2 rw;  // Corner radius, border width
	uniform vec4 colors[3]; // Fillcolor, FrameColor, SurfaceColor
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
        // d1 is distance from frame
		float d1 = sdRoundedBox(p, halfbox, rr);

		// hb2 is the inside of the frame.
		vec2 hb2 = vec2(halfbox.x-bw, halfbox.y-bw);
		float d2 = sdRoundedBox(p, hb2, rr-bw);
		if (d1>-0.5) {
			vec4 col =  colors[2]; // vec4(0.0, 0.0, 0.0, 0.0);
			fragColor = mix(colors[1], col, clamp(d1+0.5, 0, 1));
		}
		if (d2<0.5) { 
            // We are inside box. Mix with border to smooth border
			fragColor = mix(colors[1], colors[0], clamp(0.5-d2, 0, 1));
		}
	}
	` + "\x00"

// FragShadowSource is the fragment shader used for shaddows. Used together with VertRectSource
var FragShadowSource = `
	#version 330
	in vec4 gl_FragCoord;
	out vec4 fragColor;

	uniform vec2 pos;
	uniform vec2 halfbox;
    uniform vec4 rws;  // Corner radius, border width, shaddow size, shadow alfa
	uniform vec2 resolution;
	uniform vec4 colors[2]; // Fillcolor, FrameColor

	float sdRoundedBox( in vec2 p, in vec2 b, in float r ) {
		vec2 q = abs(p)-b+r;
		return min(max(q.x,q.y),0.0) + length(max(q,0.0)) - r;
	}

	void main() {
		float sw = rws.z;  // Shadow width
        float rr = rws.x;  // Corner radius
        fragColor = colors[0];
        vec2 p = vec2(gl_FragCoord.x-pos.x, resolution.y-gl_FragCoord.y-pos.y);
		// halfbox includes shadow. hb1 subtracts shadow to get frame size.
        vec2 hb1 = vec2(halfbox.x, halfbox.y);
        // Now d1 is distance from shadow center
		float d1 = sdRoundedBox(p, hb1, rr);
		if (d1>-sw) {
            // smoothstep(0, -sw,d1) gives smooth shadow outside rectangle
	        fragColor[3] = fragColor[3]*max(0.0, smoothstep(0, -sw, d1));
		}
	}
	` + "\x00"
