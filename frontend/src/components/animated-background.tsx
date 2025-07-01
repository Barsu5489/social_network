'use client';

import React, { useRef, useEffect } from 'react';
import { useTheme } from 'next-themes';

const AnimatedBackground = () => {
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const { resolvedTheme } = useTheme();

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    let width = window.innerWidth;
    let height = window.innerHeight;
    canvas.width = width;
    canvas.height = height;

    const handleResize = () => {
      width = window.innerWidth;
      height = window.innerHeight;
      canvas.width = width;
      canvas.height = height;
    };

    window.addEventListener('resize', handleResize);

    let t = 0;
    let animationFrameId: number;

    // This is a more direct and faithful translation of the provided formula.
    const renderPoint = (x: number, y: number) => {
      // k and e are input coordinates, derived from the main loop counter i.
      // The constants are preserved from the original formula.
      const k = x / 4 - 12.5;
      const e = y / 9;
      const o = Math.hypot(k, e) / 9;
      
      // c is the angle, dependent on coordinates and time.
      const c = o / 5 + e / 4 - t / 8;
      
      // q is the magnitude/radius, with complex oscillations.
      const q_part1 = (3 / (k === 0 ? 0.1 : k)) * Math.sin(y); // Avoid division by zero
      const q_part2 = k * (1 + Math.cos(y) / 3 + Math.sin(e + o * 4 - t * 2));
      const q = x / 3 + 99 + q_part1 + q_part2;

      // Convert polar-like coordinates (q, c) to cartesian and center on screen.
      const pointX = q * Math.cos(c) + width / 2;
      const pointY = (q + 49) * Math.sin(c) * Math.cos(c) - q / 3 + 30 * o + height / 2;
      
      ctx.fillRect(pointX, pointY, 1.5, 1.5);
    }

    const draw = () => {
      // Use original animation speed
      t += Math.PI / 90; 
      
      const isDark = resolvedTheme === 'dark';
      const bgColor = isDark ? 'hsl(215 28% 17%)' : 'hsl(206 67% 98%)';
      const pointColor = isDark ? 'rgba(93, 173, 226, 0.15)' : 'rgba(46, 134, 193, 0.15)';

      ctx.fillStyle = bgColor;
      ctx.fillRect(0, 0, width, height);
      
      ctx.fillStyle = pointColor;

      // Use original loop count and input value generation
      for (let i = 20000; i--;) {
        const x = i % 100;
        const y = i / 350;
        renderPoint(x, y);
      }

      animationFrameId = requestAnimationFrame(draw);
    };

    draw();

    return () => {
      cancelAnimationFrame(animationFrameId);
      window.removeEventListener('resize', handleResize);
    };
  }, [resolvedTheme]);

  return <canvas ref={canvasRef} className="fixed top-0 left-0 w-full h-full -z-50" />;
};

export default AnimatedBackground;
