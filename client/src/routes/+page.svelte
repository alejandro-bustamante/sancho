<script>
  import NotificationList from '$lib/components/NotificationList.svelte';
  import SearchBar from '$lib/components/SearchBar.svelte';
  import TrackList from '$lib/components/TrackList.svelte';
</script>

<div class="bg-dark-gray text-light min-h-screen flex flex-col items-center pt-5 px-4">
  <NotificationList />
  <SearchBar />
  <TrackList />
</div>

