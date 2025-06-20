<!-- src/routes/+layout.svelte -->
    <script>
    // You can add global stores, authentication logic, etc. here
    import { onMount } from 'svelte';
    
    let mounted = false;
    
    onMount(() => {
      mounted = true;
    });
  </script>
  
  <svelte:head>
    <title>SkyClone - Social Media</title>
    <meta name="description" content="A modern social media platform" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
  </svelte:head>
  
  <main class="app">
    {#if mounted}
      <slot />
    {/if}
  </main>
  
  <style>
    :global(html) {
      --primary-blue: #1185fe;
      --primary-blue-hover: #0d6efd;
      --bg-primary: #ffffff;
      --bg-secondary: #f8f9fa;
      --text-primary: #0f1419;
      --text-secondary: #536471;
      --border-color: #e1e8ed;
      --border-hover: #cfd9de;
      --shadow-light: 0 1px 3px rgba(0, 0, 0, 0.12);
      --shadow-medium: 0 4px 6px rgba(0, 0, 0, 0.1);
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
      line-height: 1.4;
    }
  
    :global(body) {
      margin: 0;
      padding: 0;
      background-color: var(--bg-primary);
      color: var(--text-primary);
      overflow-x: hidden;
    }
  
    :global(*) {
      box-sizing: border-box;
    }
  
    :global(button) {
      cursor: pointer;
      border: none;
      font-family: inherit;
      font-weight: 600;
      border-radius: 24px;
      transition: all 0.2s ease;
    }
  
    :global(input) {
      font-family: inherit;
      border: 1px solid var(--border-color);
      border-radius: 6px;
      padding: 12px 16px;
      font-size: 17px;
      transition: border-color 0.2s ease;
    }
  
    :global(input:focus) {
      outline: none;
      border-color: var(--primary-blue);
      box-shadow: 0 0 0 2px rgba(17, 133, 254, 0.2);
    }
  
    .app {
      min-height: 100vh;
      display: flex;
      flex-direction: column;
    }
  
    @media (prefers-color-scheme: dark) {
      :global(html) {
        --bg-primary: #000000;
        --bg-secondary: #16181c;
        --text-primary: #e7e9ea;
        --text-secondary: #71767b;
        --border-color: #2f3336;
        --border-hover: #3d4144;
      }
    }
  </style>