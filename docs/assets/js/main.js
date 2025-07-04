// Termonaut Homepage JavaScript
// Interactive features and animations

document.addEventListener('DOMContentLoaded', function() {
    // Initialize all components
    initParticles();
    initTypedAnimation();
    initDemoTerminal();
    initGitHubStats();
    initCopyButtons();
    initScrollAnimations();
    initMobileMenu();
    
    console.log('🚀 Termonaut homepage loaded successfully!');
});

// Particles.js configuration
function initParticles() {
    if (typeof particlesJS !== 'undefined') {
        particlesJS('particles-js', {
            particles: {
                number: {
                    value: 50,
                    density: {
                        enable: true,
                        value_area: 800
                    }
                },
                color: {
                    value: '#58a6ff'
                },
                shape: {
                    type: 'circle',
                    stroke: {
                        width: 0,
                        color: '#000000'
                    }
                },
                opacity: {
                    value: 0.3,
                    random: false,
                    anim: {
                        enable: false,
                        speed: 1,
                        opacity_min: 0.1,
                        sync: false
                    }
                },
                size: {
                    value: 2,
                    random: true,
                    anim: {
                        enable: false,
                        speed: 40,
                        size_min: 0.1,
                        sync: false
                    }
                },
                line_linked: {
                    enable: true,
                    distance: 150,
                    color: '#58a6ff',
                    opacity: 0.2,
                    width: 1
                },
                move: {
                    enable: true,
                    speed: 1,
                    direction: 'none',
                    random: false,
                    straight: false,
                    out_mode: 'out',
                    bounce: false,
                    attract: {
                        enable: false,
                        rotateX: 600,
                        rotateY: 1200
                    }
                }
            },
            interactivity: {
                detect_on: 'canvas',
                events: {
                    onhover: {
                        enable: true,
                        mode: 'grab'
                    },
                    onclick: {
                        enable: true,
                        mode: 'push'
                    },
                    resize: true
                },
                modes: {
                    grab: {
                        distance: 140,
                        line_linked: {
                            opacity: 0.5
                        }
                    },
                    push: {
                        particles_nb: 4
                    }
                }
            },
            retina_detect: true
        });
    }
}

// Typed.js animation for hero terminal
function initTypedAnimation() {
    if (typeof Typed !== 'undefined') {
        const commands = [
            'termonaut setup^1000',
            'termonaut tui^1000',
            'termonaut stats --weekly^1000',
            'git commit -m "Level up!"^1000',
            'termonaut github sync^1000'
        ];
        
        new Typed('#typed-command', {
            strings: commands,
            typeSpeed: 50,
            backSpeed: 30,
            backDelay: 2000,
            loop: true,
            showCursor: true,
            cursorChar: '▋',
            onStringTyped: function(arrayPos, self) {
                updateTerminalOutput(arrayPos);
            }
        });
    }
}

// Update terminal output based on typed command
function updateTerminalOutput(commandIndex) {
    const outputs = [
        // termonaut setup
        `<span style="color: #3fb950;">✓</span> Shell integration installed successfully!
<span style="color: #3fb950;">✓</span> Configuration file created
<span style="color: #58a6ff;">🚀</span> Welcome to Termonaut! Your terminal journey begins now.`,
        
        // termonaut tui
        `<span style="color: #58a6ff;">🎮</span> Launching interactive dashboard...
<span style="color: #f85149;">━━━━━━━━━━</span> <span style="color: #3fb950;">Level 8 Astronaut</span> <span style="color: #f85149;">━━━━━━━━━━</span>
<span style="color: #d29922;">XP:</span> 2,150 / 2,500 <span style="color: #3fb950;">████████████░░░░</span> 86%
<span style="color: #58a6ff;">Commands Today:</span> 127 | <span style="color: #d29922;">Streak:</span> 12 days 🔥`,
        
        // termonaut stats --weekly
        `<span style="color: #58a6ff;">📊 Weekly Terminal Stats</span>
─────────────────────────────────
Commands: 1,247 <span style="color: #3fb950;">⬆ +23%</span>
Active Time: 28h 15m
Sessions: 42
New Commands: 15 <span style="color: #d29922;">⭐</span>
Top: git (234), ls (156), cd (98)`,
        
        // git commit
        `[main 7f8a9b2] Level up!
 3 files changed, 47 insertions(+), 12 deletions(-)
<span style="color: #3fb950;">🎉 Achievement Unlocked: Git Commander!</span>
<span style="color: #58a6ff;">+50 XP</span> | You're now a Level 9 Space Commander! 🚀`,
        
        // termonaut github sync
        `<span style="color: #58a6ff;">🔄</span> Syncing with GitHub...
<span style="color: #3fb950;">✓</span> Profile updated
<span style="color: #3fb950;">✓</span> Badges generated
<span style="color: #3fb950;">✓</span> Stats synchronized
<span style="color: #d29922;">🏷️</span> Dynamic badges available at your repo!`
    ];
    
    const outputElement = document.getElementById('terminal-output');
    if (outputElement && outputs[commandIndex]) {
        outputElement.innerHTML = outputs[commandIndex];
    }
}

// Demo terminal functionality
function initDemoTerminal() {
    const demoButtons = document.querySelectorAll('.demo-btn');
    const demoTerminalBody = document.getElementById('demo-terminal-body');
    
    const demoContent = {
        stats: `<div class="terminal-line">
    <span class="prompt">$</span>
    <span class="command">termonaut stats</span>
</div>
<div class="terminal-output" style="color: #00ff00;">
🚀 Today's Terminal Stats (2024-07-04)
─────────────────────────────────────
Commands Executed: 127 🎯
Active Time: 3h 42m ⏱️
Session Count: 4 📱
New Commands: 3 ⭐
Current Streak: 12 days 🔥

Top Commands:
git (23) <span style="color: #58a6ff;">████████████████████████████████</span>
ls (18)  <span style="color: #58a6ff;">███████████████████████████</span>
cd (15)  <span style="color: #58a6ff;">██████████████████████</span>
vim (12) <span style="color: #58a6ff;">████████████████</span>

🎮 Level 8 Astronaut (2,150 XP)
Progress to Level 9: <span style="color: #3fb950;">████████████░░░░</span> 75%
</div>`,
        
        tui: `<div class="terminal-line">
    <span class="prompt">$</span>
    <span class="command">termonaut tui</span>
</div>
<div class="terminal-output" style="color: #00ff00;">
<span style="color: #f85149;">╭─────────────────────────────────────────────╮</span>
<span style="color: #f85149;">│</span>  <span style="color: #58a6ff;">🚀 Termonaut Dashboard</span>                    <span style="color: #f85149;">│</span>
<span style="color: #f85149;">├─────────────────────────────────────────────┤</span>
<span style="color: #f85149;">│</span>  <span style="color: #3fb950;">Level 8 Astronaut</span> | <span style="color: #d29922;">XP: 2,150/2,500</span>     <span style="color: #f85149;">│</span>
<span style="color: #f85149;">│</span>  <span style="color: #3fb950;">████████████░░░░</span> 86%                   <span style="color: #f85149;">│</span>
<span style="color: #f85149;">├─────────────────────────────────────────────┤</span>
<span style="color: #f85149;">│</span>  📊 Stats  🏆 Achievements  🎭 Easter Eggs  <span style="color: #f85149;">│</span>
<span style="color: #f85149;">╰─────────────────────────────────────────────╯</span>

<span style="color: #58a6ff;">Use arrow keys to navigate, Enter to select</span>
</div>`,
        
        achievements: `<div class="terminal-line">
    <span class="prompt">$</span>
    <span class="command">termonaut achievements</span>
</div>
<div class="terminal-output" style="color: #00ff00;">
🏆 Achievement Gallery (12/20 unlocked)

<span style="color: #3fb950;">✓ 🚀 First Launch</span>        Welcome aboard!
<span style="color: #3fb950;">✓ 🌟 Explorer</span>           Used 50 unique commands
<span style="color: #3fb950;">✓ 🏆 Century</span>            100 commands in one day
<span style="color: #3fb950;">✓ 🔥 Streak Keeper</span>      7-day usage streak
<span style="color: #3fb950;">✓ 👨‍🚀 Space Commander</span>   Reached level 10
<span style="color: #3fb950;">✓ 🧬 Git Commander</span>      Master git operations

<span style="color: #8b949e;">⏳ 🪐 Cosmic Explorer</span>    30-day usage streak
<span style="color: #8b949e;">⏳ ⚡ Lightning Fast</span>     500 commands in one day
<span style="color: #8b949e;">⏳ 🛸 Master Navigator</span>   Reach level 25

<span style="color: #d29922;">Progress: 60% complete</span>
</div>`,
        
        'easter-eggs': `<div class="terminal-line">
    <span class="prompt">$</span>
    <span class="command">git push origin main</span>
</div>
<div class="terminal-output" style="color: #00ff00;">
Enumerating objects: 15, done.
Counting objects: 100% (15/15), done.
Delta compression using up to 8 threads
Compressing objects: 100% (8/8), done.
Writing objects: 100% (9/9), 1.23 KiB | 1.23 MiB/s, done.
Total 9 (delta 6), reused 0 (delta 0), pack-reused 0
To github.com:username/project.git
   7f8a9b2..a1b2c3d  main -> main

<span style="color: #f85149;">🎭 Easter Egg Triggered!</span>
<span style="color: #d29922;">🚀 "Houston, we have a deployment!" 
   Another successful mission to production orbit!</span>
<span style="color: #58a6ff;">+25 XP Bonus</span> for deployment mastery!
</div>`
    };
    
    demoButtons.forEach(button => {
        button.addEventListener('click', function() {
            // Remove active class from all buttons
            demoButtons.forEach(btn => btn.classList.remove('active'));
            // Add active class to clicked button
            this.classList.add('active');
            
            // Update demo content
            const demoType = this.getAttribute('data-demo');
            if (demoTerminalBody && demoContent[demoType]) {
                demoTerminalBody.innerHTML = demoContent[demoType];
            }
        });
    });
}

// Fetch GitHub stats
async function initGitHubStats() {
    try {
        const response = await fetch('https://api.github.com/repos/oiahoon/termonaut');
        const data = await response.json();
        
        // Update stats
        document.getElementById('stars-count').textContent = data.stargazers_count || '0';
        document.getElementById('forks-count').textContent = data.forks_count || '0';
        document.getElementById('issues-count').textContent = data.open_issues_count || '0';
        
        // Update hero stats
        const heroStars = document.getElementById('github-stars');
        if (heroStars) {
            heroStars.textContent = `⭐ ${data.stargazers_count || '0'}`;
        }
        
        // Fetch releases count
        const releasesResponse = await fetch('https://api.github.com/repos/oiahoon/termonaut/releases');
        const releasesData = await releasesResponse.json();
        document.getElementById('releases-count').textContent = releasesData.length || '0';
        
    } catch (error) {
        console.log('GitHub API rate limit or network error:', error);
        // Set fallback values
        document.getElementById('stars-count').textContent = '50+';
        document.getElementById('forks-count').textContent = '10+';
        document.getElementById('releases-count').textContent = '15+';
        document.getElementById('issues-count').textContent = '5';
        
        const heroStars = document.getElementById('github-stars');
        if (heroStars) {
            heroStars.textContent = '⭐ 50+';
        }
    }
}

// Copy to clipboard functionality
function initCopyButtons() {
    const copyButtons = document.querySelectorAll('.copy-btn, .copy-btn-small');
    
    copyButtons.forEach(button => {
        button.addEventListener('click', async function() {
            const textToCopy = this.getAttribute('data-copy');
            
            try {
                await navigator.clipboard.writeText(textToCopy);
                
                // Visual feedback
                const originalContent = this.innerHTML;
                this.innerHTML = this.classList.contains('copy-btn-small') ? '✓' : 
                    '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20,6 9,17 4,12"></polyline></svg>';
                this.style.color = '#3fb950';
                
                setTimeout(() => {
                    this.innerHTML = originalContent;
                    this.style.color = '';
                }, 2000);
                
            } catch (err) {
                console.error('Failed to copy text: ', err);
                
                // Fallback for older browsers
                const textArea = document.createElement('textarea');
                textArea.value = textToCopy;
                document.body.appendChild(textArea);
                textArea.select();
                document.execCommand('copy');
                document.body.removeChild(textArea);
                
                // Visual feedback
                const originalContent = this.innerHTML;
                this.innerHTML = this.classList.contains('copy-btn-small') ? '✓' : 
                    '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20,6 9,17 4,12"></polyline></svg>';
                this.style.color = '#3fb950';
                
                setTimeout(() => {
                    this.innerHTML = originalContent;
                    this.style.color = '';
                }, 2000);
            }
        });
    });
}

// Scroll animations
function initScrollAnimations() {
    const observerOptions = {
        threshold: 0.1,
        rootMargin: '0px 0px -50px 0px'
    };
    
    const observer = new IntersectionObserver(function(entries) {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add('animated');
            }
        });
    }, observerOptions);
    
    // Add animation classes to elements
    const animatedElements = document.querySelectorAll('.feature-card, .stat-card, .install-method, .step');
    animatedElements.forEach(el => {
        el.classList.add('animate-on-scroll');
        observer.observe(el);
    });
}

// Mobile menu functionality
function initMobileMenu() {
    const navToggle = document.querySelector('.nav-toggle');
    const navMenu = document.querySelector('.nav-menu');
    
    if (navToggle && navMenu) {
        navToggle.addEventListener('click', function() {
            navMenu.classList.toggle('active');
            navToggle.classList.toggle('active');
        });
        
        // Close menu when clicking on links
        const navLinks = document.querySelectorAll('.nav-link');
        navLinks.forEach(link => {
            link.addEventListener('click', function() {
                navMenu.classList.remove('active');
                navToggle.classList.remove('active');
            });
        });
    }
}

// Smooth scrolling for anchor links
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault();
        const target = document.querySelector(this.getAttribute('href'));
        if (target) {
            target.scrollIntoView({
                behavior: 'smooth',
                block: 'start'
            });
        }
    });
});

// Navbar scroll effect
window.addEventListener('scroll', function() {
    const navbar = document.querySelector('.navbar');
    if (window.scrollY > 100) {
        navbar.style.background = 'rgba(13, 17, 23, 0.98)';
        navbar.style.backdropFilter = 'blur(20px)';
    } else {
        navbar.style.background = 'rgba(13, 17, 23, 0.95)';
        navbar.style.backdropFilter = 'blur(10px)';
    }
});

// Easter egg: Konami code
let konamiCode = [];
const konamiSequence = [
    'ArrowUp', 'ArrowUp', 'ArrowDown', 'ArrowDown',
    'ArrowLeft', 'ArrowRight', 'ArrowLeft', 'ArrowRight',
    'KeyB', 'KeyA'
];

document.addEventListener('keydown', function(e) {
    konamiCode.push(e.code);
    
    if (konamiCode.length > konamiSequence.length) {
        konamiCode.shift();
    }
    
    if (konamiCode.join(',') === konamiSequence.join(',')) {
        // Easter egg activated!
        showEasterEgg();
        konamiCode = [];
    }
});

function showEasterEgg() {
    const easterEgg = document.createElement('div');
    easterEgg.innerHTML = `
        <div style="
            position: fixed;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            background: var(--bg-tertiary);
            border: 2px solid var(--accent-primary);
            border-radius: var(--radius-lg);
            padding: var(--spacing-xl);
            text-align: center;
            z-index: 10000;
            box-shadow: var(--shadow-glow);
            animation: fadeInUp 0.5s ease;
        ">
            <div style="font-size: 3rem; margin-bottom: var(--spacing-md);">🎉</div>
            <h3 style="color: var(--accent-primary); margin-bottom: var(--spacing-md);">Easter Egg Unlocked!</h3>
            <p style="color: var(--text-secondary); margin-bottom: var(--spacing-lg);">
                You found the secret Konami code! 🕹️<br>
                True terminal ninjas know the classics.
            </p>
            <button onclick="this.parentElement.parentElement.remove()" style="
                background: var(--accent-primary);
                color: var(--bg-primary);
                border: none;
                padding: var(--spacing-sm) var(--spacing-lg);
                border-radius: var(--radius-md);
                cursor: pointer;
                font-weight: 600;
            ">Awesome! 🚀</button>
        </div>
        <div style="
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background: rgba(0, 0, 0, 0.8);
            z-index: 9999;
        " onclick="this.parentElement.remove()"></div>
    `;
    
    document.body.appendChild(easterEgg);
    
    // Auto-remove after 10 seconds
    setTimeout(() => {
        if (easterEgg.parentElement) {
            easterEgg.remove();
        }
    }, 10000);
}

// Performance monitoring
if ('performance' in window) {
    window.addEventListener('load', function() {
        setTimeout(function() {
            const perfData = performance.getEntriesByType('navigation')[0];
            console.log(`🚀 Page loaded in ${Math.round(perfData.loadEventEnd - perfData.loadEventStart)}ms`);
        }, 0);
    });
}

// Service worker registration (for future PWA features)
if ('serviceWorker' in navigator) {
    window.addEventListener('load', function() {
        // Uncomment when service worker is ready
        // navigator.serviceWorker.register('/sw.js')
        //     .then(registration => console.log('SW registered'))
        //     .catch(error => console.log('SW registration failed'));
    });
}
