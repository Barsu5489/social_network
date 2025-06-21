<script>
  // You can add props or animations here later
</script>

<div class="hero-container">
  <div class="overlay">
    <h1>Welcome to Skiep</h1>
    <p>
    The plartorm that gives you freedom to express your inner self what you were 
    meant to be. Post the most exiting moments of life, we only have this one 😁.
    </p>
  </div>

  <!-- Geometric animated background -->
  <div class="geometric-bg">
    {#each Array(10) as _, i}
      <div class="shape" style="--i: {i}"></div>
    {/each}
  </div>
</div>

