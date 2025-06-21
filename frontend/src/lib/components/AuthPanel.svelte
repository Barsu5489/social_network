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
