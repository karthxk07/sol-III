#!/bin/bash

# Setup script for activating virtual environment and running video processor
# Usage: ./setup.sh [optional: path/to/video.mp4]

VENV_PATH="./.venv"
PYTHON_SCRIPT="video_processor.py"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Video Processor Setup ===${NC}"

# Check if virtual environment exists
if [ ! -d "$VENV_PATH" ]; then
  echo -e "${YELLOW}Virtual environment not found. Creating...${NC}"
  python3 -m venv "$VENV_PATH"
fi

# Activate virtual environment
echo -e "${GREEN}Activating virtual environment...${NC}"
source "$VENV_PATH/bin/activate"

# Verify activation
if [ "$VIRTUAL_ENV" != "" ]; then
  echo -e "${GREEN}✓ Virtual environment activated: $VIRTUAL_ENV${NC}"
else
  echo -e "${RED}✗ Failed to activate virtual environment${NC}"
  exit 1
fi

# Check if requirements.txt exists
if [ ! -f "requirements.txt" ]; then
  echo -e "${RED}✗ requirements.txt not found${NC}"
  exit 1
fi

# Check and install dependencies from requirements.txt
echo -e "${GREEN}Checking dependencies from requirements.txt...${NC}"

# Check if all requirements are satisfied
pip install -q --dry-run -r requirements.txt 2>/dev/null

if [ $? -ne 0 ]; then
  echo -e "${YELLOW}Installing missing dependencies...${NC}"
  pip install -r requirements.txt
  if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ All dependencies installed successfully${NC}"
  else
    echo -e "${RED}✗ Failed to install dependencies${NC}"
    exit 1
  fi
else
  echo -e "${GREEN}✓ All dependencies already satisfied${NC}"
fi

# Optional: Run the Python script if argument provided
if [ "$1" != "" ]; then
  echo -e "${GREEN}Running video processor on: $1${NC}"
  python "$PYTHON_SCRIPT" "$1"
else
  echo -e "${GREEN}Setup complete! Run your Python script with:${NC}"
  echo -e "${YELLOW}python $PYTHON_SCRIPT${NC}"
fi

# Keep shell in venv (if sourced)
echo -e "${GREEN}Virtual environment is active. Type 'deactivate' to exit.${NC}"
