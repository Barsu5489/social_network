<script>
  // You can add props or animations here later
</script>

<div class="hero-container">
  <div class="overlay">
    <h1>Welcome to website</h1>
    <p>
      Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam
      nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat.
    </p>
  </div>

  <!-- Geometric animated background -->
  <div class="geometric-bg">
    {#each Array(10) as _, i}
      <div class="shape" style="--i: {i}"></div>
    {/each}
  </div>
</div>

<style>
  .hero-container {
    height: 100%;
    width: 100%;
    background: linear-gradient(135deg, #6366f1, #ec4899);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 4rem 2rem;
    position: relative;
    color: white;
  }

  .overlay {
    z-index: 10;
    max-width: 480px;
    text-align: left;
  }

  h1 {
    font-size: 2.5rem;
    font-weight: 700;
    margin-bottom: 1rem;
  }

  p {
    font-size: 1.125rem;
    opacity: 0.85;
  }

  .geometric-bg {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    overflow: hidden;
    z-index: 1;
  }

  .shape {
    position: absolute;
    width: 60px;
    height: 8px;
    background-color: #f97316;
    opacity: 0.3;
    transform: rotate(45deg);
    animation: float 20s linear infinite;
    animation-delay: calc(-2s * var(--i));
    top: calc(10% * var(--i));
    left: calc(10% * var(--i));
  }

  @keyframes float {
    0% {
      transform: translateY(0) rotate(45deg);
    }
    100% {
      transform: translateY(-200%) rotate(45deg);
    }
  }
</style>
