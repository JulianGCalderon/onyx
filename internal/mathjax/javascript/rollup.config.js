import commonjs from '@rollup/plugin-commonjs';
import { nodeResolve } from '@rollup/plugin-node-resolve';


export default {
  input: "main.js",
  output: {
    file: 'bundle.js',
    name: 'mathjax',
    format: 'iife'
  },
  plugins: [commonjs(), nodeResolve()]
};

