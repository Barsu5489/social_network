<script>
  import { login } from '../services/auth/authService';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';

  let email = '';
  let password = '';
  let rememberMe = false;
  let error = '';
  let loading = false;

  async function handleLogin() {
    if (!email || !password) {
      error = 'Please fill in all fields';
      return;
    }

    if (password.length < 6) {
      error = 'Password must be at least 6 characters';
      return;
    }

    loading = true;
    error = '';

    try {
      await login({ email, password });
      goto('/');
    } catch (err) {
      error = err.message || 'Login failed';
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    // Check if user is already authenticated
    // This is just a placeholder - you'll need to implement this
    // if (isAuthenticated()) {
    //   goto('/');
    // }
  });
</script>

<style>
  .auth-panel {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    padding: 2rem;
    max-width: 400px;
    margin: 0 auto;
  }

  .auth-header {
    font-size: 1.25rem;
    font-weight: 600;
    color: #6366f1;
    text-transform: uppercase;
    margin-bottom: 2rem;
    text-align: center;
  }

  .form-control {
    position: relative;
    margin-bottom: 1.5rem;
    width: 100%;
  }

  .form-control input {
    width: 100%;
    padding: 0.75rem 1rem;
    padding-left: 3rem;
    border: 1px solid #e5e7eb;
    border-radius: 0.375rem;
    font-size: 1rem;
    transition: border-color 0.2s ease;
  }

  .form-control input:focus {
    border-color: #6366f1;
    outline: none;
  }

  .form-control svg {
    position: absolute;
    top: 50%;
    left: 1rem;
    transform: translateY(-50%);
    color: #6b7280;
  }

  .remember-me {
    display: flex;
    align-items: center;
    margin-bottom: 1rem;
  }

  .remember-me input {
    margin-right: 0.5rem;
  }

  .forgot-password {
    text-align: right;
    margin-bottom: 1.5rem;
  }

  .forgot-password a {
    color: #6b7280;
    text-decoration: none;
    transition: color 0.2s ease;
  }

  .forgot-password a:hover {
    color: #6366f1;
  }

  .login-button {
    width: 100%;
    padding: 0.75rem;
    background: linear-gradient(90deg, #6366f1 0%, #4338ca 100%);
    color: white;
    border: none;
    border-radius: 0.375rem;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.2s ease, transform 0.2s ease;
  }

  .login-button:hover {
    background: linear-gradient(90deg, #4338ca 0%, #6366f1 100%);
    transform: scale(1.02);
  }

  .login-button:disabled {
    background: #a5b4fc;
    cursor: not-allowed;
  }

  .error {
    color: red;
    margin-bottom: 1rem;
    text-align: center;
  }
</style>

<div class="auth-panel">
  <h2 class="auth-header">USER LOGIN</h2>

  {#if error}
    <div class="error">{error}</div>
  {/if}

  <form on:submit|preventDefault={handleLogin}>
    <div class="form-control">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-6 h-6">
        <path fill-rule="evenodd" d="M18 10a8 8 0 100-16 8 8 0 000 16zm-6-3a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0zm-1 5.25a2.25 2.25 0 100-4.5 2.25 2.25 0 000 4.5z" clip-rule="evenodd" />
      </svg>
      <input
        type="email"
        placeholder="Enter username"
        bind:value={email}
      />
    </div>

    <div class="form-control">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-6 h-6">
        <path fill-rule="evenodd" d="M12 1.5a5.25 5.25 0 00-5.25 5.25v9a5.25 5.25 0 1010.5 0v-9a5.25 5.25 0 00-5.25-5.25zm3 8.25a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z" clip-rule="evenodd" />
      </svg>
      <input
        type="password"
        placeholder="Enter password"
        bind:value={password}
      />
    </div>

    <div class="remember-me">
      <input type="checkbox" bind:checked={rememberMe} />
      <label>Remember me</label>
    </div>

    <div class="forgot-password">
      <a href="#">Forgot password?</a>
    </div>

    <button
      type="submit"
      class="login-button"
      disabled={loading}
    >
      {#if loading}
        <span>Loading...</span>
      {:else}
        <span>LOGIN</span>
      {/if}
    </button>
  </form>
</div>
