const esbuild = require("esbuild");

esbuild.build({
    entryPoints: ["src/main.ts"],
    bundle: true,
    outfile: "dist/bundle.js",
    format: "esm",
    target: "es2015",
    sourcemap: true,
    minify: false,
}).then(() => {
    console.log("Build successful!");
}).catch(() => {
    process.exit(1);
});