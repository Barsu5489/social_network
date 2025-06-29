<!-- Layout.svelte -->
<script>
  import Navbar from '$lib/components/Navbar.svelte';
  import LeftSidebar from '$lib/components/LeftSidebar.svelte';
  import { user, isAuthenticated, getCurrentUser, getAllPosts } from '$lib/services/api';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';

  let userData = null;
  let authStatus = false;

  user.subscribe(value => {
    userData = value;
  });
  isAuthenticated.subscribe(value => {
    authStatus = value;
  });

  onMount(async () => {
    if (!authStatus) {
      try {
        // Fetch current user profile to verify authentication
        await getCurrentUser();
      } catch (error) {
        console.error('Failed to fetch user:', error);
        goto('/login');
        return;
      }
    }

    try {
      // Fetch posts to populate feeds
      await getAllPosts();
    } catch (error) {
      console.error('Failed to load posts:', error);
    }
  });
</script>

{#if userData && authStatus}
  <div class="app-layout">
    <Navbar user={userData} />
    
    <div class="main-container">
      <LeftSidebar user={userData} />
      
      <main class="main-content">
        <slot />
      </main>
      
      <RightSidebar />
    </div>
  </div>
{:else}
  <p>Loading...</p>
{/if}
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