<script>
  import { goto } from '$app/navigation';
  import { isAuthenticated } from '$lib/services/auth/authService';
  import { onMount } from 'svelte';

  onMount(() => {
    isAuthenticated.subscribe(value => {
      if (value) {
        goto('/home');
      }
    });
  });
</script>

<svelte:head>
  <title>ConnectSphere - Connect, Share, Belong</title>
</svelte:head>

<main class="landing-page">
  <nav class="navbar">
    <div class="nav-container">
      <h1 class="logo">ConnectSphere</h1>
      <div class="nav-buttons">
        <button on:click={() => goto('/login')} class="nav-btn login-btn">Login</button>
        <button on:click={() => goto('/register')} class="nav-btn register-btn">Sign Up</button>
      </div>
    </div>
  </nav>

  <section class="hero">
    <div class="hero-content">
      <div class="hero-text">
        <h1>Connect with friends and the world around you</h1>
        <p>Share moments, discover communities, and stay connected with the people who matter most to you.</p>
        <div class="hero-buttons">
          <button on:click={() => goto('/register')} class="cta-btn primary">Get Started</button>
          <button on:click={() => goto('/login')} class="cta-btn secondary">Sign In</button>
        </div>
      </div>
      <div class="hero-visual">
        <div class="feature-grid">
          <div class="feature-card">
            <div class="feature-icon">👥</div>
            <h3>Connect</h3>
            <p>Follow friends and discover new people</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">📝</div>
            <h3>Share</h3>
            <p>Post updates, photos, and moments</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">👥</div>
            <h3>Groups</h3>
            <p>Join communities and create events</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">💬</div>
            <h3>Chat</h3>
            <p>Real-time messaging with friends</p>
          </div>
        </div>
      </div>
    </div>
  </section>

  <section class="features">
    <div class="container">
      <h2>Why Choose ConnectSphere?</h2>
      <div class="features-grid">
        <div class="feature">
          <div class="feature-icon large">🔒</div>
          <h3>Privacy First</h3>
          <p>Control who sees your content with public and private profiles</p>
        </div>
        <div class="feature">
          <div class="feature-icon large">⚡</div>
          <h3>Real-time Updates</h3>
          <p>Instant notifications and live messaging</p>
        </div>
        <div class="feature">
          <div class="feature-icon large">🌐</div>
          <h3>Rich Media</h3>
          <p>Share photos, GIFs, and create engaging posts</p>
        </div>
      </div>
    </div>
  </section>

  <footer class="footer">
    <div class="container">
      <p>&copy; 2025 ConnectSphere. Designed by swapo and brav</p>
    </div>
  </footer>
</main>

<style>
  :global(body) {
    margin: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
    background-color: #ffffff;
    color: #1f2937;
    line-height: 1.6;
  }

  .landing-page {
    min-height: 100vh;
  }

  .navbar {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border-bottom: 1px solid #e5e7eb;
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    z-index: 100;
  }

  .nav-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 1rem 2rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .logo {
    font-size: 1.75rem;
    font-weight: 700;
    background: linear-gradient(135deg, #6366f1, #ec4899);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    margin: 0;
  }

  .nav-buttons {
    display: flex;
    gap: 1rem;
  }

  .nav-btn {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 0.5rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .login-btn {
    background: transparent;
    color: #6b7280;
  }

  .login-btn:hover {
    background: #f3f4f6;
    color: #374151;
  }

  .register-btn {
    background: linear-gradient(135deg, #6366f1, #ec4899);
    color: white;
  }

  .register-btn:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3);
  }

  .hero {
    padding: 8rem 2rem 4rem;
    background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
  }

  .hero-content {
    max-width: 1200px;
    margin: 0 auto;
    display: grid;
    grid-template-columns: 1fr 1fr;;
    gap: 4rem;
    align-items: center;
  }

  .hero-text h1 {
    font-size: 3rem;
    font-weight: 800;
    color: #1f2937;
    margin-bottom: 1.5rem;
    line-height: 1.2;
  }

  .hero-text p {
    font-size: 1.25rem;
    color: #6b7280;
    margin-bottom: 2rem;
  }

  .hero-buttons {
    display: flex;
    gap: 1rem;
    flex-wrap: wrap;
  }

  .cta-btn {
    padding: 1rem 2rem;
    border: none;
    border-radius: 0.75rem;
    font-size: 1.1rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
  }

  .cta-btn.primary {
    background: linear-gradient(135deg, #6366f1, #ec4899);
    color: white;
  }

  .cta-btn.primary:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 25px rgba(99, 102, 241, 0.4);
  }

  .cta-btn.secondary {
    background: white;
    color: #6366f1;
    border: 2px solid #6366f1;
  }

  .cta-btn.secondary:hover {
    background: #6366f1;
    color: white;
  }

  .feature-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 1.5rem;
  }

  .feature-card {
    background: white;
    padding: 2rem;
    border-radius: 1rem;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
    text-align: center;
    transition: transform 0.2s ease;
  }

  .feature-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
  }

  .feature-icon {
    font-size: 2rem;
    margin-bottom: 1rem;
  }

  .feature-icon.large {
    font-size: 3rem;
  }

  .feature-card h3 {
    font-size: 1.25rem;
    font-weight: 600;
    margin-bottom: 0.5rem;
    color: #1f2937;
  }

  .feature-card p {
    color: #6b7280;
    font-size: 0.9rem;
  }

  .features {
    padding: 4rem 2rem;
    background: white;
  }

  .container {
    max-width: 1200px;
    margin: 0 auto;
  }

  .features h2 {
    text-align: center;
    font-size: 2.5rem;
    font-weight: 700;
    color: #1f2937;
    margin-bottom: 3rem;
  }

  .features-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 2rem;
  }

  .feature {
    text-align: center;
    padding: 2rem;
  }

  .feature h3 {
    font-size: 1.5rem;
    font-weight: 600;
    margin: 1rem 0;
    color: #1f2937;
  }

  .feature p {
    color: #6b7280;
    font-size: 1rem;
  }

  .footer {
    background: #1f2937;
    color: #9ca3af;
    text-align: center;
    padding: 2rem;
  }

  .footer p {
    margin: 0;
  }

  @media (max-width: 768px) {
    .hero-content {
      grid-template-columns: 1fr;
      gap: 2rem;
      text-align: center;
    }

    .hero-text h1 {
      font-size: 2rem;
    }

    .feature-grid {
      grid-template-columns: 1fr;
    }

    .nav-container {
      padding: 1rem;
    }

    .logo {
      font-size: 1.5rem;
    }

    .hero {
      padding: 6rem 1rem 2rem;
    }

    .hero-buttons {
      justify-content: center;
    }
  }
</style>
