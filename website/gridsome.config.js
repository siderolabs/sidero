// This is where project configuration and plugin options are located.
// Learn more: https://gridsome.org/docs/config

// Changes here require a server restart.
// To restart press CTRL + C in terminal and run `gridsome develop`

module.exports = {
  siteName: "Sidero",
  icon: {
    favicon: "./src/assets/favicon.png",
    touchicon: "./src/assets/favicon.png",
  },
  siteUrl: process.env.SITE_URL ? process.env.SITE_URL : "https://talos.dev",
  settings: {
    title: "Bare Metal Kubernetes",
    description: "A bare metal provisioning system for managing Kubernetes clusters",
    web: process.env.URL_WEB || false,
    twitter: "https://twitter.com/talossystems",
    github: "https://github.com/talos-systems/sidero",
    nav: {
      links: [{ path: "", title: "Docs" }],
    },
    dropdownOptions: [
      {
        version: "v0.5",
        url: "/docs/v0.5/",
        latest: false,
        prerelease: true,
      },
      {
        version: "v0.4",
        url: "/docs/v0.4/",
        latest: true,
        prerelease: false,
      },
      {
        version: "v0.3",
        url: "/docs/v0.3/",
        latest: false,
        prerelease: false,
      },
      {
        version: "v0.2",
        url: "/docs/v0.2/",
        latest: false,
        prerelease: false,
      },
      {
        version: "v0.1",
        url: "/docs/v0.1/",
        latest: false,
        prerelease: false,
      },
    ],
  },

  // Allow '.' in slugs (e.g. /docs/v0.1).
  permalinks: {
    slugify: {
      use: "slugify",
      options: { lower: true },
    },
  },

  plugins: [
    {
      use: "gridsome-source-docs",
      options: {
        baseDir: "./content/docs",
        path: "**/*.md",
        typeName: "MarkdownPage",
        pathPrefix: "/docs",
        sidebarOrder: {
          "v0.5": [
            { title: "Overview", method: "weighted" },
            { title: "Getting Started", method: "weighted" },
            { title: "Resource Configuration", method: "alphabetical" },
            { title: "Guides", method: "alphabetical" },
          ],
          "v0.4": [
            { title: "Overview", method: "weighted" },
            { title: "Getting Started", method: "weighted" },
            { title: "Resource Configuration", method: "alphabetical" },
            { title: "Guides", method: "alphabetical" },
          ],
          "v0.3": [
            { title: "Overview", method: "weighted" },
            { title: "Getting Started", method: "weighted" },
            { title: "Resource Configuration", method: "alphabetical" },
            { title: "Guides", method: "alphabetical" },
          ],
          "v0.2": [
            { title: "Overview", method: "weighted" },
            { title: "Getting Started", method: "weighted" },
            { title: "Resource Configuration", method: "alphabetical" },
            { title: "Guides", method: "alphabetical" },
          ],
          "v0.1": [
            { title: "Overview", method: "weighted" },
            { title: "Resource Configuration", method: "alphabetical" },
            { title: "Guides", method: "alphabetical" },
          ],
        },
        remark: {
          externalLinksTarget: "_blank",
          externalLinksRel: ["noopener", "noreferrer"],
          plugins: [
            "gridsome-plugin-remark-mermaid",
            "@gridsome/remark-prismjs"],
        },
      },
    },

    {
      use: "gridsome-plugin-tailwindcss",
      options: {
        tailwindConfig: "./tailwind.config.js",
      },
    },

    {
      use: "@gridsome/plugin-google-analytics",
      options: {
        id: process.env.GA_ID ? process.env.GA_ID : "XX-999999999-9",
      },
    },

    {
      use: "@gridsome/plugin-sitemap",
      options: {},
    },
  ],
};
