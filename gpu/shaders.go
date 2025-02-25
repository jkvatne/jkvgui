package gpu

var (
	RectangleFragShaderSource = `
	#version 400
	in  vec2 aRadWidth;
	in  vec4 aRect;
	in  vec4 BorderColor;
	in  vec4 FillColor;
	layout(origin_upper_left) in vec4 gl_FragCoord;
	out vec4 fragColor;
	
	// b.x = half width, b.y = half height
	float sdRoundedBox( in vec2 p, in vec2 b, in float r ) {
		vec2 q = abs(p)-b+r;
		return min(max(q.x,q.y),0.0) + length(max(q,0.0)) - r;
	}

	void main() {
		fragColor = FillColor;
		vec2 halfbox = vec2((aRect[2]-aRect[0])/2, (aRect[3]-aRect[1])/2);
		vec2 p = gl_FragCoord.xy;
		p = p-vec2((aRect[2]+aRect[0])/2, (aRect[3]+aRect[1])/2);
		float d1 = sdRoundedBox(p, halfbox, aRadWidth[0]);
		float w = aRadWidth[1];
		vec2 halfbox2 = vec2(halfbox.x-w*2, halfbox.y-2*w);
		float d2 = sdRoundedBox(p, halfbox2, aRadWidth[0]-w);
		if (d1>0.0) {
			discard;
		}
		if (d2<=0) {
			fragColor = BorderColor;
		}
	}
	` + "\x00"

	RectangleVertShaderSource = `
	#version 400
	layout(location = 1) in vec2 inPos;
	layout(location = 2) in vec2 inColorIndex;
	layout(location = 3) in vec2 inRadWidth;
	layout(location = 4) in vec4 inRect;
	out  vec2 aRadWidth;
	out  vec4 aRect;
	out  vec4 BorderColor;
	out  vec4 FillColor;
	uniform vec2 resolution;
	uniform vec4 colors[8];
	void main() {
		vec2 zeroToOne = inPos / resolution;
		vec2 zeroToTwo = zeroToOne * 2.0;
		vec2 clipSpace = zeroToTwo - 1.0;
		gl_Position = vec4(clipSpace * vec2(1, -1), 0, 1);
		BorderColor =  colors[int(inColorIndex[0])];
		FillColor =  colors[int(inColorIndex[1])];
		aRadWidth = inRadWidth;
		aRect = inRect;
	}
	` + "\x00"

	fragmentFontShader = `#version 400
	in vec2 fragTexCoord;
	out vec4 outputColor;
	uniform sampler2D tex;
	uniform vec4 textColor;
	void main()
	{    
		vec4 sampled = vec4(1.0, 1.0, 1.0, texture(tex, fragTexCoord).r);
		outputColor = textColor * sampled;
	}	
	` + "\x00"

	vertexFontShader = `#version 400
	in vec2 vert;
	in vec2 vertTexCoord;
	uniform vec2 resolution;
	out vec2 fragTexCoord;
	
	void main() {
	   vec2 zeroToOne = vert / resolution;
	   vec2 zeroToTwo = zeroToOne * 2.0;
	   vec2 clipSpace = zeroToTwo - 1.0;
	   fragTexCoord = vertTexCoord;
	   gl_Position = vec4(clipSpace * vec2(1, -1), 0, 1);
	}
	` + "\x00"
)

// vec2 halfbox = vec2((aRect[2]-aRect[0])/2, (aRect[3]-aRect[1])/2);
// vec2 p = gl_FragCoord.xy;
// p = p-vec2((aRect[2]+aRect[0])/2, (aRect[3]+aRect[1])/2);
// float d1 = sdRoundedBox(p, halfbox, aRadWidth[0]);
// float w = aRadWidth[1];
// vec2 halfbox2 = vec2(halfbox.x-w*2, halfbox.y-2*w);
// float d2 = sdRoundedBox(p, halfbox2, aRadWidth[0]-w);

var (
	RectFragShaderSource = `
	#version 400

	uniform vec2 pos;
	uniform vec2 halfbox;
    uniform vec2 rw;
	uniform vec4 colors[2];
	uniform vec2 resolution;

	//layout(origin_upper_left) 
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
	#version 400
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
