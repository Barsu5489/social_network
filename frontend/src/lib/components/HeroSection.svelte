<script>
  // You can add props or animations here later
</script>

<div class="hero-container">
  <div class="overlay">
    <h1>Welcome to Skiep</h1>
    <p>
      Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam
      nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat.
    </p>
  </div>

  <!-- Geometric animated background -->
  <div class="geometric-bg">
    {#each Array(10) as _, i}
      <div class="shape" style="--i: {i}"></div>
    {/each}
  </div>
</div>

