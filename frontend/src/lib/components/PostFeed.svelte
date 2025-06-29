<!-- PostFeed.svelte -->
<script>
  import { posts, getAllPosts } from '$lib/services/api';
  import PostCard from './PostCard.svelte';
  import { onMount } from 'svelte';

  export let user = null;

  let errorMessage = '';

  onMount(async () => {
    try {
      await getAllPosts();
    } catch (error) {
      errorMessage = 'Failed to load posts. Please try again.';
      console.error('Failed to load posts:', error);
    }
  });
</script>

<div class="post-feed">
  {#if errorMessage}
    <div class="error-message">{errorMessage}</div>
  {/if}
  
  {#each $posts as post}
    <PostCard {post} {user} />
  {/each}
</div>

<style>
  .post-feed {
    max-width: 600px;
    margin: 0 auto;
  }
  .error-message {
    color: red;
    font-size: 0.9em;
    margin-bottom: 1em;
  }
</style>