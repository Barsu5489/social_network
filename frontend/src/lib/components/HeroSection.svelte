<script>
  import { onMount } from 'svelte';

  let gradientAngle = 0;

  function updateGradient() {
    gradientAngle = (gradientAngle + 0.5) % 360;
    requestAnimationFrame(updateGradient);
  }

  onMount(() => {
    updateGradient();
  });
</script>

<style>
  :global(:root) {
    --gradient-angle: v-bind(gradientAngle + 'deg');
  }

  .hero {
    position: relative;
    width: 100%;
    height: 100%;
    background: linear-gradient(
      var(--gradient-angle, 45deg),
      #6366f1,
      #ec4899
    );
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: flex-start;
    padding: 4rem 6rem;
    color: #ffffff;
    overflow: hidden;
  }

  .hero h1 {
    font-size: 4rem;
    font-weight: 800;
    margin-bottom: 1rem;
    opacity: 0.9;
    animation: fadeInLeft 1s ease-out;
  }

  .hero p {
    font-size: 1.25rem;
    max-width: 600px;
    opacity: 0.8;
    animation: fadeInLeft 1.5s ease-out;
  }

  @keyframes fadeInLeft {
    from {
      opacity: 0;
      transform: translateX(-20px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }

  .geometric-elements {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;
  }

  .geometric-element {
    position: absolute;
    background: #f97316;
    opacity: 0.8;
    animation: rotate 20s linear infinite, pulse 3s ease-in-out infinite;
  }

  .geometric-element:nth-child(2) {
    background: #ec4899;
    animation-delay: 5s;
  }

  .geometric-element:nth-child(3) {
    background: #f97316;
    animation-delay: 10s;
  }

  @keyframes rotate {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }

  @keyframes pulse {
    0%, 100% {
      transform: scale(1);
    }
    50% {
      transform: scale(1.1);
    }
  }
</style>

<div class="hero">
  <div class="geometric-elements">
    <div class="geometric-element" style="width: 100px; height: 100px; top: 10%; left: 15%; transform: rotate(20deg);"></div>
    <div class="geometric-element" style="width: 150px; height: 50px; top: 30%; left: 5%; transform: rotate(-30deg);"></div>
    <div class="geometric-element" style="width: 80px; height: 80px; top: 60%; left: 70%; transform: rotate(45deg);"></div>
  </div>
  <h1>Welcome to website</h1>
  <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
</div>
