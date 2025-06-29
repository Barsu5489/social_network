
<script>
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { getPost } from '../../lib/services/postService';
  import PostCard from '../../lib/components/PostCard.svelte';

  let post = null;

  onMount(async () => {
    const postId = $page.params.id;
    try {
      post = await getPost(postId);
    } catch (error) {
      console.error('Error fetching post:', error);
    }
  });
</script>

<div class="post-page-container">
  {#if post}
    <PostCard {post} />
  {:else}
    <p>Loading post...</p>
  {/if}
</div>

<style>
  .post-page-container {
    max-width: 800px;
    margin: 80px auto;
    padding: 20px;
  }
</style>
