<script>

import Navbar from '$lib/components/Navbar.svelte';
import LeftSidebar from '$lib/components/LeftSidebar.svelte';

  
  export let user = {
    name: 'Swapo El Nenjez',
    email: 'swapo@example.com'
  };
</script>

<div class="app-layout">
  <Navbar {user} />
  
  <div class="main-container">
    <LeftSidebar {user} />
    
    <main class="main-content">
      <slot />
    </main>
    
    <RightSidebar />
  </div>
</div>

<style>
  .app-layout {
    min-height: 100vh;
    background: #f0f2f5;
  }
  
  .main-container {
    display: flex;
    position: relative;
    padding-top: 56px; /* Height of navbar */
  }
  
  .main-content {
    flex: 1;
    margin-left: 280px; /* Width of left sidebar */
    margin-right: 280px; /* Width of right sidebar */
    padding: 16px;
    max-width: calc(100vw - 560px);
    min-height: calc(100vh - 56px);
  }
  
  /* Responsive design */
  @media (max-width: 1200px) {
    .main-content {
      margin-right: 0;
      max-width: calc(100vw - 280px);
    }
  }
  
  @media (max-width: 900px) {
    .main-content {
      margin-left: 0;
      margin-right: 0;
      max-width: 100vw;
      padding: 16px 8px;
    }
  }
</style>