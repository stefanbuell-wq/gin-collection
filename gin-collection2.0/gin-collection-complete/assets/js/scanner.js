// Barcode Scanner using Quagga2
class BarcodeScanner {
    constructor() {
        this.isScanning = false;
        this.lastResult = null;
        this.lastResultTime = 0;
    }

    start() {
        if (this.isScanning) return;

        const container = document.getElementById('scanner-container');
        
        Quagga.init({
            inputStream: {
                name: "Live",
                type: "LiveStream",
                target: container,
                constraints: {
                    width: { min: 640 },
                    height: { min: 480 },
                    facingMode: "environment",
                    aspectRatio: { min: 1, max: 2 }
                }
            },
            decoder: {
                readers: [
                    "ean_reader",
                    "ean_8_reader",
                    "code_128_reader",
                    "code_39_reader",
                    "upc_reader",
                    "upc_e_reader"
                ],
                debug: {
                    drawBoundingBox: true,
                    showFrequency: true,
                    drawScanline: true,
                    showPattern: true
                }
            },
            locator: {
                patchSize: "medium",
                halfSample: true
            },
            numOfWorkers: 4,
            frequency: 10,
            locate: true
        }, (err) => {
            if (err) {
                console.error('Scanner initialization error:', err);
                alert('Kamera konnte nicht gestartet werden: ' + err.message);
                return;
            }
            
            console.log('Scanner initialized');
            Quagga.start();
            this.isScanning = true;
        });

        Quagga.onDetected((result) => {
            const code = result.codeResult.code;
            const now = Date.now();
            
            // Avoid duplicate scans within 2 seconds
            if (code === this.lastResult && now - this.lastResultTime < 2000) {
                return;
            }

            this.lastResult = code;
            this.lastResultTime = now;

            // Play success sound (optional)
            this.playBeep();

            // Visual feedback
            container.style.border = '5px solid #27ae60';
            setTimeout(() => {
                container.style.border = 'none';
            }, 500);

            console.log('Barcode detected:', code);
            
            // Send to app
            if (window.app) {
                this.stop();
                window.app.handleBarcodeScan(code);
            }
        });

        Quagga.onProcessed((result) => {
            const drawingCtx = Quagga.canvas.ctx.overlay;
            const drawingCanvas = Quagga.canvas.dom.overlay;

            if (result) {
                if (result.boxes) {
                    drawingCtx.clearRect(0, 0, drawingCanvas.width, drawingCanvas.height);
                    
                    // Draw detection boxes
                    result.boxes.filter(box => box !== result.box).forEach(box => {
                        Quagga.ImageDebug.drawPath(box, { x: 0, y: 1 }, drawingCtx, {
                            color: "green",
                            lineWidth: 2
                        });
                    });
                }

                if (result.box) {
                    Quagga.ImageDebug.drawPath(result.box, { x: 0, y: 1 }, drawingCtx, {
                        color: "#00F",
                        lineWidth: 2
                    });
                }

                if (result.codeResult && result.codeResult.code) {
                    Quagga.ImageDebug.drawPath(result.line, { x: 'x', y: 'y' }, drawingCtx, {
                        color: 'red',
                        lineWidth: 3
                    });
                }
            }
        });
    }

    stop() {
        if (!this.isScanning) return;

        Quagga.stop();
        this.isScanning = false;
        
        // Clear the container
        const container = document.getElementById('scanner-container');
        container.innerHTML = '';
    }

    playBeep() {
        // Create a simple beep sound
        const audioContext = new (window.AudioContext || window.webkitAudioContext)();
        const oscillator = audioContext.createOscillator();
        const gainNode = audioContext.createGain();

        oscillator.connect(gainNode);
        gainNode.connect(audioContext.destination);

        oscillator.frequency.value = 800;
        oscillator.type = 'sine';

        gainNode.gain.setValueAtTime(0.3, audioContext.currentTime);
        gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + 0.1);

        oscillator.start(audioContext.currentTime);
        oscillator.stop(audioContext.currentTime + 0.1);
    }
}

// Make scanner globally available
window.BarcodeScanner = new BarcodeScanner();
