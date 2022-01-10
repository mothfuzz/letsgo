#version 100

attribute vec3 position;
attribute vec2 texcoord;

varying vec2 frag_texcoord;

uniform mat4 MVP;

void main() {
    gl_Position = MVP * vec4(position, 1.0);
    frag_texcoord = texcoord;
}
