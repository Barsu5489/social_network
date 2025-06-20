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
    background-color: #ffffff;
    color: #6b7280;
  }

  .landing-page {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
    position: relative;
  }

  .hero-section {
    flex: 0 0 60%;
    position: relative;
    overflow: hidden;
  }

  .auth-panel {
    flex: 0 0 40%;
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
    position: absolute;
    bottom: 0;
    width: 100%;
  }

  @media (max-width: 1024px) {
    .landing-page {
      flex-direction: column;
    }

    .hero-section,
    .auth-panel {
      flex: 1 1 auto;
      width: 100%;
    }
  }
</style>
