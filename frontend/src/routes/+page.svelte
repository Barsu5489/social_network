<script>
  import HeroSection from '$lib/components/HeroSection.svelte';
  import AuthPanel from '$lib/components/AuthPanel.svelte';
</script>

<main class="landing-page">
  <section class="hero-section">
    <HeroSection />
  </section>

  <section class="auth-panel">
    <AuthPanel />
  </section>

  <footer class="credit">
    <p>designed by swapo</p>
  </footer>
</main>

<style>
  :global(body) {
    margin: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, 'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol', 'Noto Color Emoji';
    background-color: #752a2a;
    color: #6b7280;
  }

  .landing-page {
    display: grid;
    place-items: center;
    min-height: 100vh;
    position: relative;
  }

  .hero-section {
    position: relative;
    overflow: hidden;
  }

  
  .auth-panel {
    background: #ffffff;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    padding: 4rem 2rem;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  .credit {
    text-align: center;
    padding: 1rem;
    background: #f8fafc;
    color: #6b7280;
    font-size: 0.875rem;
    width: 100%;
    position: absolute;
    bottom: 0;
  }

  @media (max-width: 1024px) {
    .landing-page {
      display: flex;
      flex-direction: column;
    }

    .hero-section,
    .auth-panel {
      width: 100%;
    }
  }
</style>
