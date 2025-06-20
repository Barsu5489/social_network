<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { isAuthenticated } from '$lib/services/auth/authService';

  let authenticated = false;

  $: if ($isAuthenticated) {
    authenticated = true;
  }

  onMount(() => {
    if (!authenticated) {
      goto('/login');
    }
  });
</script>

{#if authenticated}
<main>
  <h1>Your Profile</h1>
  <p>This is your profile page. You can view and edit your information here.</p>
</main>
{/if}
