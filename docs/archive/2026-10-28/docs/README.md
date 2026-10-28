# 🚀 Termonaut Homepage

This directory contains the source code for the Termonaut project homepage - a modern, geek-style website that showcases the features and capabilities of Termonaut.

## 🎨 Design Philosophy

The homepage follows a **terminal-inspired geek aesthetic** with:

- **Dark theme** with terminal-like colors and fonts
- **Interactive terminal simulations** showing real Termonaut usage
- **Animated elements** and smooth transitions
- **Responsive design** that works on all devices
- **Performance-optimized** with minimal dependencies

## 🏗️ Architecture

### File Structure
```
docs/
├── index.html              # Main HTML structure
├── assets/
│   ├── css/
│   │   └── style.css      # Main stylesheet (terminal theme)
│   ├── js/
│   │   └── main.js        # Interactive functionality
│   └── images/            # Static assets (logos, screenshots)
├── README.md              # This file
└── sitemap.xml            # Generated during build
```

### Key Features

#### 🎬 Interactive Elements
- **Typed.js animations** in hero terminal
- **Live demo terminal** with switchable content
- **Particle.js background** for visual appeal
- **Smooth scroll animations** and hover effects

#### 📊 Dynamic Content
- **Real-time GitHub stats** fetched via API
- **Copy-to-clipboard** functionality for code snippets
- **Responsive terminal layouts** that adapt to screen size
- **Easter eggs** including Konami code activation

#### 🎯 Performance Features
- **Lazy loading** for non-critical resources
- **Optimized animations** with CSS transforms
- **Minimal JavaScript** with vanilla JS (no heavy frameworks)
- **Compressed assets** and efficient caching

## 🚀 Deployment

The homepage is automatically deployed to GitHub Pages using GitHub Actions:

### Deployment Workflow
1. **Trigger**: Push to `main` branch or changes to `docs/` directory
2. **Build**: Process HTML, optimize assets, generate sitemap
3. **Test**: Validate HTML structure and check links
4. **Deploy**: Upload to GitHub Pages
5. **Monitor**: Run Lighthouse performance tests

### GitHub Actions Features
- ✅ **Automatic deployment** on code changes
- ✅ **HTML validation** and link checking
- ✅ **Asset optimization** (CSS minification)
- ✅ **Performance monitoring** with Lighthouse CI
- ✅ **Dynamic content generation** (stats, timestamps)

## 🎨 Customization

### Color Scheme
The homepage uses CSS custom properties for easy theming:

```css
:root {
  --bg-primary: #0d1117;      /* Main background */
  --text-primary: #f0f6fc;    /* Primary text */
  --accent-primary: #58a6ff;  /* Brand accent */
  --terminal-bg: #0c0c0c;     /* Terminal background */
  --terminal-text: #00ff00;   /* Terminal text */
}
```

### Typography
- **Headings**: Inter (modern sans-serif)
- **Code/Terminal**: JetBrains Mono (monospace)
- **Body text**: Inter with optimized line-height

### Responsive Breakpoints
- **Desktop**: 1024px+
- **Tablet**: 768px - 1023px
- **Mobile**: < 768px

## 🔧 Development

### Local Development
```bash
# Serve locally (Python)
cd docs
python -m http.server 8000

# Or with Node.js
npx serve .

# Visit: http://localhost:8000
```

### Testing
```bash
# HTML validation
npx html-validate index.html

# Lighthouse testing
npx lighthouse http://localhost:8000 --output html

# Accessibility testing
npx axe-cli http://localhost:8000
```

### Adding New Sections
1. Add HTML structure to `index.html`
2. Add corresponding styles to `assets/css/style.css`
3. Add interactive functionality to `assets/js/main.js`
4. Update navigation links if needed

## 📊 Performance Targets

The homepage is optimized for:

- **Performance**: > 90 Lighthouse score
- **Accessibility**: > 95 Lighthouse score
- **SEO**: > 95 Lighthouse score
- **Load Time**: < 2 seconds on 3G
- **Bundle Size**: < 500KB total

## 🎭 Interactive Features

### Terminal Simulations
- **Hero Terminal**: Animated typing with realistic command outputs
- **Demo Terminal**: Switchable content showing different Termonaut features
- **Code Blocks**: Copy-to-clipboard functionality

### Animations
- **Scroll Animations**: Elements fade in as they enter viewport
- **Hover Effects**: Cards lift and glow on hover
- **Loading States**: Smooth transitions for dynamic content

### Easter Eggs
- **Konami Code**: Classic cheat code triggers special message
- **Performance Monitoring**: Console logs for debugging
- **GitHub API Integration**: Live stats with fallback values

## 🌐 SEO & Social

### Meta Tags
- Complete Open Graph tags for social sharing
- Twitter Card support
- Structured data for search engines
- Optimized meta descriptions

### Social Media
- **Open Graph Image**: Custom designed for social sharing
- **Twitter Cards**: Rich previews on Twitter
- **LinkedIn**: Professional sharing optimization

## 🔒 Security & Privacy

- **No tracking scripts** or analytics by default
- **Minimal external dependencies** (only CDN fonts and libraries)
- **Content Security Policy** headers via GitHub Pages
- **HTTPS only** with automatic redirects

## 🚀 Future Enhancements

### Planned Features
- [ ] **PWA Support**: Service worker for offline functionality
- [ ] **Dark/Light Toggle**: User preference switching
- [ ] **Interactive Tutorials**: Step-by-step Termonaut guides
- [ ] **Community Showcase**: User-generated content section
- [ ] **Multilingual Support**: i18n for global audience

### Performance Improvements
- [ ] **Image Optimization**: WebP format with fallbacks
- [ ] **Critical CSS**: Inline above-the-fold styles
- [ ] **Preload Hints**: Optimize resource loading priority
- [ ] **Service Worker**: Advanced caching strategies

## 📞 Support

For homepage-related issues:

1. **Bug Reports**: [GitHub Issues](https://github.com/oiahoon/termonaut/issues)
2. **Feature Requests**: [GitHub Discussions](https://github.com/oiahoon/termonaut/discussions)
3. **Documentation**: [Project Wiki](https://github.com/oiahoon/termonaut/wiki)

## 📄 License

The homepage code is released under the same MIT License as the main Termonaut project.

---

**Built with ❤️ for the terminal community** 🚀
