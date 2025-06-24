<script>
  import { goto } from '$app/navigation';
  import { login } from '$lib/services/auth/authService';
  
  let email = "";
  let password = "";
  let isLoading = false;
  let error = "";

  const handleLogin = async () => {
    if (!email || !password) {
      error = "Please fill in all fields";
      return;
    }

    isLoading = true;
    error = "";

    try {
      await login({ email, password });
      // Handle successful login - redirect to profile
      goto('/profile');
    } catch (err) {
      error = err.message || 'Login failed. Please try again.';
    } finally {
      isLoading = false;
    }
  };
</script>

<div class="login-container">
  <div class="login-card">
    <div class="login-header">
      <h1>Welcome Back</h1>
      <p>Sign in to your ConnectSphere account</p>
    </div>

    <form on:submit|preventDefault={handleLogin} class="login-form">
      {#if error}
        <div class="error-message">
          {error}
        </div>
      {/if}

      <div class="form-group">
        <label for="email">Email Address</label>
        <input 
          id="email"
          type="email" 
          bind:value={email}
          placeholder="Enter your email"
          required
          disabled={isLoading}
        />
      </div>

      <div class="form-group">
        <label for="password">Password</label>
        <input 
          id="password"
          type="password" 
          bind:value={password}
          placeholder="Enter your password"
          required
          minlength="6"
          disabled={isLoading}
        />
      </div>

      <button type="submit" class="login-btn" disabled={isLoading}>
        {#if isLoading}
          <span class="spinner"></span>
          Signing in...
        {:else}
          Sign In
        {/if}
      </button>
    </form>

    <div class="login-footer">
      <p>Don't have an account? <a href="/register">Sign up here</a></p>
      <a href="/" class="back-link">← Back to home</a>
    </div>
  </div>
</div>

<style>
  .login-container {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    padding: 2rem;
  }

  .login-card {
    background: white;
    border-radius: 1rem;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
    padding: 3rem;
    width: 100%;
    max-width: 420px;
  }

  .login-header {
    text-align: center;
    margin-bottom: 2rem;
  }

  .login-header h1 {
    font-size: 2rem;
    font-weight: 700;
    color: #1f2937;
    margin-bottom: 0.5rem;
  }

  .login-header p {
    color: #6b7280;
    margin: 0;
  }

  .login-form {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .error-message {
    background: #fef2f2;
    border: 1px solid #fecaca;
    color: #dc2626;
    padding: 0.75rem;
    border-radius: 0.5rem;
    font-size: 0.875rem;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .form-group label {
    font-weight: 600;
    color: #374151;
    font-size: 0.875rem;
  }

  .form-group input {
    padding: 0.875rem;
    border: 2px solid #e5e7eb;
    border-radius: 0.5rem;
    font-size: 1rem;
    transition: all 0.2s ease;
  }

  .form-group input:focus {
    outline: none;
    border-color: #6366f1;
    box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
  }

  .form-group input:disabled {
    background-color: #f9fafb;
    cursor: not-allowed;
  }

  .login-btn {
    background: linear-gradient(135deg, #6366f1, #ec4899);
    color: white;
    border: none;
    padding: 1rem;
    border-radius: 0.5rem;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
  }

  .login-btn:hover:not(:disabled) {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(99, 102, 241, 0.4);
  }

  .login-btn:disabled {
    opacity: 0.7;
    cursor: not-allowed;
    transform: none;
  }

  .spinner {
    width: 1rem;
    height: 1rem;
    border: 2px solid transparent;
    border-top: 2px solid currentColor;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .login-footer {
    text-align: center;
    margin-top: 2rem;
    padding-top: 2rem;
    border-top: 1px solid #e5e7eb;
  }

  .login-footer p {
    color: #6b7280;
    margin-bottom: 1rem;
  }

  .login-footer a {
    color: #6366f1;
    text-decoration: none;
    font-weight: 600;
  }

  .login-footer a:hover {
    text-decoration: underline;
  }

  .back-link {
    font-size: 0.875rem;
    color: #9ca3af;
  }

  @media (max-width: 480px) {
    .login-container {
      padding: 1rem;
    }

    .login-card {
      padding: 2rem 1.5rem;
    }

    .login-header h1 {
      font-size: 1.75rem;
    }
  }
</style>
