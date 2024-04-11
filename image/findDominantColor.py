import cv2
import os
import cairosvg
import numpy as np
from sklearn.cluster import KMeans

# read all the images from the ./images/icons folder
image_dir = './images/icons'
total = 0
totalExtracted = 0

def convert_svg_to_png(svg_file, png_file):
    """Convert SVG to PNG using CairoSVG."""
    try:
        cairosvg.svg2png(url=svg_file, write_to=png_file)
        return True
    except Exception as e:
        print("Error converting SVG to PNG:", e)
        return False

def convert_jpeg_to_png(jpeg_file, png_file):
    """Convert JPEG/JPG to PNG using OpenCV."""
    try:
        image = cv2.imread(jpeg_file)
        cv2.imwrite(png_file, image)
        return True
    except Exception as e:
        print("Error converting JPEG/JPG to PNG:", e)
        return False

def extract_dominant_color(image_path, k=3):
    """Extract dominant color from an image using KMeans."""
    # Load the image
    image = cv2.imread(image_path)
    if image is None:
        print(f"Error: Unable to load image from {image_path}")
        return None
    
    image = cv2.cvtColor(image, cv2.COLOR_BGR2RGB)  # Convert to RGB

    # Reshape the image to be a list of pixels
    pixels = image.reshape(-1, 3)

    # Apply KMeans to find the dominant colors
    kmeans = KMeans(n_clusters=k)
    kmeans.fit(pixels)

    # Get the dominant colors
    dominant_colors = kmeans.cluster_centers_.astype(int)
    return dominant_colors

for image_name in os.listdir(image_dir):
    image_path = os.path.join(image_dir, image_name)
    total += 1
    
    # Define SVG and PNG file paths
    svg_file = image_path
    png_file = os.path.splitext(image_path)[0] + ".png"

    # Convert SVG to PNG
    if svg_file.lower().endswith('.svg'):
        if convert_svg_to_png(svg_file, png_file):
            print("SVG converted to PNG successfully!")
            dominant_colors_svg = extract_dominant_color(png_file)
            if dominant_colors_svg is not None:
                print("Dominant colors from SVG:", dominant_colors_svg)
        else:
            print("Failed to convert SVG to PNG.")
    
    # Convert JPEG/JPG to PNG
    elif image_path.lower().endswith('.jpeg') or image_path.lower().endswith('.jpg'):
        if convert_jpeg_to_png(image_path, png_file):
            print("JPEG/JPG converted to PNG successfully!")
            dominant_colors_jpeg = extract_dominant_color(png_file)
            if dominant_colors_jpeg is not None:
                print("Dominant colors from JPEG/JPG:", dominant_colors_jpeg)
        else:
            print("Failed to convert JPEG/JPG to PNG.")
    
    # Extract dominant colors directly for PNG files
    elif image_path.lower().endswith('.png'):
        dominant_colors = extract_dominant_color(image_path)
        if dominant_colors is not None:
            totalExtracted += 1
            print("Dominant colors of", image_name, ":", dominant_colors)
    else:
        print(f"Unsupported file format for {image_name}")

print("Total:", total)
print("Extracted:", totalExtracted)
