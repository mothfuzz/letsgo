#version 100
precision highp float;

varying vec2 frag_texcoord;

uniform sampler2D tex;

void main() {
    gl_FragColor = texture2D(tex, frag_texcoord);
    if(gl_FragColor.a == 0.0) {
        discard;
    }
}
