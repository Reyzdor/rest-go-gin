class OptimizedRegisterSlider {
    constructor() {
        this.slider = document.getElementById('registerSlider');
        this.currentSlide = 0;
        this.intervalTime = 6000; 
        this.slideInterval = null;
        this.isTransitioning = false;
        
        this.photos = [
            './static/media/photos/1.jpg',
            './static/media/photos/11.jpg',
            './static/media/photos/15.jpg',
            './static/media/photos/8.jpg',
            './static/media/photos/6.jpg',
            './static/media/photos/9.jpg'
        ];
        
        this.preloadAllImages().then(() => {
            this.init();
        }).catch(() => {
            this.init();
        });
    }
    
  
    async preloadAllImages() {
        const promises = this.photos.map(photo => {
            return new Promise((resolve, reject) => {
                const img = new Image();
                img.onload = () => {
                    resolve();
                };
                img.onerror = () => {
                    resolve(); 
                };
                img.src = photo;
            });
        });
        
        await Promise.all(promises);
    }
    
   
    init() {
        if (!this.slider) return;
        
        this.slider.style.minHeight = '100vh';
        this.slider.style.minWidth = '100%';
        
        this.createOptimizedSlides();
        
        this.showSlide(0, false);
        
        setTimeout(() => {
            this.startAutoSlide();
        }, 100);
        
        this.addOptimizedEventListeners();
        
        if (window.gc) {
            window.gc();
        }
    }
    
  
    createOptimizedSlides() {
        this.slider.innerHTML = '';
        
        this.slidesContainer = document.createElement('div');
        this.slidesContainer.className = 'slides-container-optimized';
        
        this.realSlidesCount = this.photos.length;
        
        Object.assign(this.slidesContainer.style, {
            display: 'flex',
            width: `${(this.realSlidesCount + 1) * 100}%`,
            height: '100%',
            transition: 'transform 0.6s cubic-bezier(0.25, 0.46, 0.45, 0.94)',
            willChange: 'transform', 
            backfaceVisibility: 'hidden', 
            perspective: '1000px'
        });
        
        this.slides = [];
        
        for (let i = 0; i <= this.realSlidesCount; i++) {
            const photoIndex = i % this.realSlidesCount;
            const slide = document.createElement('div');
            slide.className = `slide-optimized slide-${i}`;
            
            Object.assign(slide.style, {
                width: `${100 / (this.realSlidesCount + 1)}%`,
                height: '100%',
                backgroundImage: `url(${this.photos[photoIndex]})`,
                backgroundSize: 'cover',
                backgroundPosition: 'center',
                backgroundRepeat: 'no-repeat',
                flexShrink: '0',
                willChange: 'transform', 
                backfaceVisibility: 'hidden',
                transform: 'translate3d(0,0,0)' 
            });
            
            this.preloadBackground(slide);
            
            this.slidesContainer.appendChild(slide);
            this.slides.push(slide);
        }
        
        this.slider.appendChild(this.slidesContainer);
        this.totalSlides = this.slides.length;
    }
    
  
    preloadBackground(slide) {
        const bgImg = new Image();
        const bgStyle = slide.style.backgroundImage;
        const url = bgStyle.slice(5, -2); 
        bgImg.src = url;
    }
    
    /**
     * @param {number} index 
     * @param {boolean} animate 
     */
    showSlide(index, animate = true) {
        if (this.isTransitioning || index === this.currentSlide) return;
        
        this.isTransitioning = true;
        
        const translateX = -((index * 100) / this.totalSlides);
        
        if (animate) {
            this.slidesContainer.style.transform = `translate3d(${translateX}%, 0, 0)`;
            
            let startTime = null;
            const duration = 600;
            
            const animateTransition = (timestamp) => {
                if (!startTime) startTime = timestamp;
                const progress = timestamp - startTime;
                
                if (progress < duration) {
                    requestAnimationFrame(animateTransition);
                } else {
                    this.handleTransitionEnd(index);
                }
            };
            
            requestAnimationFrame(animateTransition);
        } else {
            this.slidesContainer.style.transition = 'none';
            this.slidesContainer.style.transform = `translate3d(${translateX}%, 0, 0)`;
            
            setTimeout(() => {
                this.currentSlide = index;
                this.isTransitioning = false;
                this.slidesContainer.style.transition = 'transform 0.6s cubic-bezier(0.25, 0.46, 0.45, 0.94)';
            }, 10);
        }
    }
    
   
    handleTransitionEnd(index) {
        this.currentSlide = index;
        
        if (index === this.realSlidesCount) {
            this.slidesContainer.style.transition = 'none';
            this.slidesContainer.style.transform = 'translate3d(0, 0, 0)';
            
            requestAnimationFrame(() => {
                requestAnimationFrame(() => {
                    this.slidesContainer.style.transition = 'transform 0.6s cubic-bezier(0.25, 0.46, 0.45, 0.94)';
                    this.currentSlide = 0;
                    this.isTransitioning = false;
                });
            });
        } else {
            this.isTransitioning = false;
        }
    }
    

    nextSlide() {
        if (this.isTransitioning) return;
        
        let nextIndex = this.currentSlide + 1;
        
        if (nextIndex > this.realSlidesCount) {
            nextIndex = 0;
        }
        
        this.showSlide(nextIndex);
    }

    startAutoSlide() {
        this.stopAutoSlide();
        
        let lastTime = Date.now();
        
        const slideLoop = () => {
            const currentTime = Date.now();
            const delta = currentTime - lastTime;
            
            if (delta >= this.intervalTime && !this.isTransitioning) {
                this.nextSlide();
                lastTime = currentTime;
            }
            
            this.slideInterval = requestAnimationFrame(slideLoop);
        };
        
        this.slideInterval = requestAnimationFrame(slideLoop);
    }
    
 
    stopAutoSlide() {
        if (this.slideInterval) {
            cancelAnimationFrame(this.slideInterval);
            this.slideInterval = null;
        }
    }
    

    addOptimizedEventListeners() {
        const options = { passive: true };
        
        this.slider.addEventListener('mouseenter', () => {
            this.stopAutoSlide();
        }, options);
        
        this.slider.addEventListener('mouseleave', () => {
            if (!this.isTransitioning) {
                this.startAutoSlide();
            }
        }, options);
        
        document.addEventListener('visibilitychange', () => {
            if (document.hidden) {
                this.stopAutoSlide();
            } else {
                if (!this.isTransitioning) {
                    this.startAutoSlide();
                }
            }
        }, options);
        
        this.slidesContainer.addEventListener('transitionend', () => {
            if (!this.isTransitioning) {
                this.slidesContainer.style.willChange = 'auto';
                this.slides.forEach(slide => {
                    slide.style.willChange = 'auto';
                });
            }
        }, options);
        
        this.slider.addEventListener('touchstart', () => {
            this.stopAutoSlide();
        }, options);
        
        this.slider.addEventListener('touchend', () => {
            setTimeout(() => {
                if (!this.isTransitioning) {
                    this.startAutoSlide();
                }
            }, 1000);
        }, options);
    }
    

    destroy() {
        this.stopAutoSlide();
        this.slider.innerHTML = '';
        this.slides = null;
        this.slidesContainer = null;
    }
}

document.addEventListener('DOMContentLoaded', () => {
    if (document.readyState === 'complete') {
        window.optimizedRegisterSlider = new OptimizedRegisterSlider();
    } else {
        window.addEventListener('load', () => {
            window.optimizedRegisterSlider = new OptimizedRegisterSlider();
        });
    }
});