set -e #Exit if a command fails

BINARY_NAME="server"
OUTPUT_DIR="./bin"
MAIN_PATH="./cmd/api"

BINARY_PATH="${OUTPUT_DIR}/${BINARY_NAME}"
echo "🧹 Tidying Go modules..."

go mod tidy

echo "Building application"
mkdir -p "$OUTPUT_DIR"


go build -o "$BINARY_PATH" "$MAIN_PATH"

echo "✅ Build successful!"

echo "🌐 Starting server on port 8000..."

PORT=8000 "$BINARY_PATH"