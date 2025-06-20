<script>
    let isLogin = true;
    let username = "";
    let email = "";
    let password = "";
    let remember = false;
  
    const handleSubmit = () => {
      if (isLogin) {
        alert(`Logging in with ${username}, remember: ${remember}`);
      } else {
        alert(`Registering with ${username}, email: ${email}`);
      }
    };
  </script>
  
  <div class="auth-container">
    <div class="tab-toggle">
      <button on:click={() => isLogin = true} class:is-active={isLogin}>Login</button>
      <button on:click={() => isLogin = false} class:is-active={!isLogin}>Register</button>
    </div>
  
    <h2>{isLogin ? 'USER LOGIN' : 'REGISTER'}</h2>
  
    <form on:submit|preventDefault={handleSubmit}>
      <div class="input-group">
        <span class="icon">👤</span>
        <input type="text" placeholder="Enter username" bind:value={username} required />
      </div>
  
      {#if !isLogin}
        <div class="input-group">
          <span class="icon">📧</span>
          <input type="email" placeholder="Enter email" bind:value={email} required />
        </div>
      {/if}
  
      <div class="input-group">
        <span class="icon">🔒</span>
        <input type="password" placeholder="Enter password" bind:value={password} required minlength="6" />
      </div>
  
      {#if isLogin}
        <div class="options">
          <label><input type="checkbox" bind:checked={remember} /> Remember</label>
          <a href="#">Forgot password?</a>
        </div>
      {/if}
  
      <button type="submit" class="login-button">{isLogin ? 'LOGIN' : 'REGISTER'}</button>
    </form>
  </div>
  
  <style>
    .auth-container {
      width: 100%;
      max-width: 400px;
      display: flex;
      flex-direction: column;
      gap: 1.5rem;
    }
  
    .tab-toggle {
      display: flex;
      justify-content: space-around;
      margin-bottom: 0.5rem;
    }
  
    .tab-toggle button {
      flex: 1;
      padding: 0.5rem;
      border: none;
      background: transparent;
      font-weight: 600;
      color: #6b7280;
      border-bottom: 2px solid transparent;
      cursor: pointer;
    }
  
    .tab-toggle button.is-active {
      color: #6366f1;
      border-color: #6366f1;
    }
  
    h2 {
      text-align: center;
      color: #6366f1;
      font-weight: 600;
      font-size: 1.25rem;
    }
  
    form {
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }
  
    .input-group {
      display: flex;
      align-items: center;
      border: 1px solid #e5e7eb;
      padding: 0.5rem;
      border-radius: 0.375rem;
      background: #f9fafb;
      gap: 0.5rem;
    }
  
    .input-group input {
      flex: 1;
      border: none;
      outline: none;
      background: transparent;
    }
  
    .options {
      display: flex;
      justify-content: space-between;
      font-size: 0.875rem;
      color: #6b7280;
    }
  
    .login-button {
      background: linear-gradient(to right, #6366f1, #ec4899);
      color: white;
      border: none;
      padding: 0.75rem;
      font-weight: 600;
      border-radius: 0.5rem;
      cursor: pointer;
      transition: transform 0.2s ease, background 0.2s ease;
    }
  
    .login-button:hover {
      transform: scale(1.03);
    }
  
    .icon {
      font-size: 1.25rem;
    }
  </style>
  