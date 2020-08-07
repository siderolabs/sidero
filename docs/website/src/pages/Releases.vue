<template>
  <Layout>
    <div>
      <div>
        <h1>Stable</h1>
        <ul>
          <li
            v-for="item in $page.github.repository.releases.nodes"
            :key="item.name"
          >
            <a v-if="!item.isPrerelease" :href="item.url">{{ item.name }}</a>
          </li>
        </ul>
      </div>
      <div>
        <h1>Pre-Release</h1>
        <ul>
          <li
            v-for="item in $page.github.repository.releases.nodes"
            :key="item.name"
          >
            <a v-if="item.isPrerelease" :href="item.url">{{ item.name }}</a>
          </li>
        </ul>
      </div>
    </div>
  </Layout>
</template>

<script>
export default {
  metaInfo: {
    title: "Releases",
  },
};
</script>

<page-query>
query {
  github {
    repository(owner: "talos-systems", name: "sidero") {
      releases(first: 10, orderBy: {field: CREATED_AT, direction: DESC}) {
        nodes {
          name
          url
          publishedAt
          isPrerelease
        }
      }
    }
  }
}
</page-query>
