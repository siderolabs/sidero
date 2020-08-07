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
  siteUrl: process.env.SITE_URL ? process.env.SITE_URL : "https://example.com",
  settings: {
    web: process.env.URL_WEB || false,
    twitter: "https://twitter.com/talossystems",
    github: "https://github.com/talos-systems/sidero",
    nav: {
      links: [
        { path: "/docs/", title: "Docs" },
        { path: "/releases/", title: "Releases" },
      ],
    },
    sidebar: [
      {
        name: "docs",
        sections: [
          {
            title: "Getting Started",
            items: [
              "/docs/",
              "/docs/installation/",
              "/docs/architecture/",
              "/docs/concepts/",
            ],
          },
          {
            title: "Configuration",
            items: [
              "/docs/environments/",
              "/docs/servers/",
              "/docs/serverclasses/",
              "/docs/metadata/",
            ],
          },
          {
            title: "Guides",
            items: ["/docs/bootstrapping/", "/docs/patching/"],
          },
        ],
      },
    ],
  },
  plugins: [
    {
      use: "@gridsome/source-filesystem",
      options: {
        baseDir: "./content",
        path: "**/*.md",
        typeName: "MarkdownPage",
        remark: {
          externalLinksTarget: "_blank",
          externalLinksRel: ["noopener", "noreferrer"],
          plugins: ["@gridsome/remark-prismjs"],
        },
      },
    },

    {
      use: "@gridsome/source-graphql",
      options: {
        url: "https://api.github.com/graphql",
        fieldName: "github",
        typeName: "github",
        headers: {
          Authorization: `Bearer ${process.env["GITHUB_TOKEN"]}`,
        },
      },
    },

    {
      use: "gridsome-plugin-tailwindcss",
      options: {
        tailwindConfig: "./tailwind.config.js",
        purgeConfig: {
          // Prevent purging of prism classes.
          whitelistPatternsChildren: [/token$/],
        },
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
