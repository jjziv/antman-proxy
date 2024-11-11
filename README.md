# antman-proxy

An image proxy that allows users to make requests to resize an image, modify an image's quality, or change an image's format.
Antman Proxy is CDN friendly and manages its own internal cache for better performance.

## Usage

To resize an image, make a GET request to: `/resize?url=[encoded_image_url]&width=[width]&height=[height]&quality=[quality]&format=[format]`

### Supported Parameters:
- `url`: URL-encoded image URL
- `width`: Desired width in pixels (optional if height is specified)
- `height`: Desired height in pixels (optional if width is specified)
- `quality`: JPEG/WebP quality (1-100, default: 85)
- `format`: Output format (jpeg, png, webp, **default:** jpeg)

### Examples:
- Resize by width with custom quality (JPEG):
`/resize?url=https%3A%2F%2Fexample.com%2Fimage.jpg&width=800&quality=90`

- Convert to WebP format:
`/resize?url=https%3A%2F%2Fexample.com%2Fimage.jpg&width=800&format=webp`

- Convert to PNG with exact dimensions:
`/resize?url=https%3A%2F%2Fexample.com%2Fimage.jpg&width=800&height=600&format=png`

## Features
- Automatic image resizing
- Smart caching system
- Multiple output formats (JPEG, PNG, WebP)
- Flexible dimension control
- Adjustable quality settings
- Rate limiting protection
- CDN-optimized responses
- Efficient browser caching

## CDN Integration
- Optimized cache headers for CDN delivery
- ETag support for efficient caching
- CORS enabled for cross-origin requests
- Long-term caching for processed images
- Short-term caching for HTML content

## Limitations
- Maximum 60 requests per minute per IP
- Only trusted domains are allowed 
- Maximum dimensions: 2000x2000 pixels
- Quality range: 1-100
- Supported formats: JPEG, PNG, WebP