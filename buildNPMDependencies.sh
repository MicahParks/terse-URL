# Build the dependencies code.
npx webpack

# Move the precompiled code to where it can be hosted.
mv dist/main.js frontend/
