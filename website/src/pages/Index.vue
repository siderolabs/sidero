<template>
  <Layout>
    <div class="pt-8 md:pt-16">
      <div class="flex flex-col items-center">
        <div class="flex flex-col items-center mb-2 text-ui-primary">
          <Logo :width="80" />
          <h2 class="text-3xl text-6xl font-black tracking-tighter border-none">
            Sidero Metal
          </h2>
        </div>
        <h1 class="text-4xl text-center lg:text-5xl">
          Kubernetes Bare Metal Lifecycle Management.
        </h1>
        <p class="text-xl font-medium text-center">Automated. Secure. Simple.</p>
<p class="text-xl font-medium text-center" style="width: 50%;">Sidero Metal uses Cluster API to automate bare metal server provisioning and lifecycle management.
Clusters are provisioned (or re-provisioned) automatically, delivering a secure Kubernetes deployment.</p>
<p class="text-xl font-medium text-center" style="width: 50%;">
Sidero includes a metadata service, PXE and TFTP servers, as well as BMC and IPMI management for automation.</p>
        <div class="flex justify-center mt-8">
          <g-link
            :to="latestDocs()"
            class="flex items-center px-6 py-4 ml-auto text-2xl font-bold leading-none text-white border rounded-lg shadow-lg bg-ui-primary border-ui-primary transition-all duration-200 ease-out transform hover:shadow-xl hover:-translate-y-1"
          >
            Get started
            <ArrowRightCircleIcon class="ml-4" size="1x" />
          </g-link>
        </div>
      </div>

      <div
        class="pt-8 mx-auto mt-8 border-t md:mt-16 md:pt-16 border-top border-ui-border max-w-screen-sm"
      ></div>

      <div class="flex flex-wrap justify-center -mx-4">
        <div
          class="flex flex-col items-center w-full px-4 mb-8 text-center md:w-1/3"
        >
          <ZapIcon size="3x" class="mb-6 text-ui-primary" />
          <h3 class="font-bold tracking-wide uppercase text-ui-primary">
            Built on Kubernetes
          </h3>
          <p class="text-lg text-left">
            Built on Kubernetes, making deployments fast and easy.
          </p>
        </div>

        <div
          class="flex flex-col items-center w-full px-4 mb-8 text-center md:w-1/3"
        >
          <ClipboardIcon size="3x" class="mb-6 text-ui-primary" />
          <h3 class="font-bold tracking-wide uppercase text-ui-primary">
            Inventory Management
          </h3>
          <p class="text-lg text-left">
            Register your servers. Classify them. Deploy them.
          </p>
        </div>

        <div
          class="flex flex-col items-center w-full px-4 mb-8 text-center md:w-1/3"
        >
          <AwardIcon size="3x" class="mb-6 text-ui-primary" />
          <h3 class="font-bold tracking-wide uppercase text-ui-primary">
            Certified Kubernetes
          </h3>
          <p class="text-lg text-left">
            Sidero Metal installs certified Kubernetes clusters secured with Talos Linux.
          </p>
        </div>
      </div>
      <p class="text-xl font-medium text-center items-center ">Sidero Metal is an Open Source project from the Sidero Labs team.
</p>
    </div>
  </Layout>
</template>

<static-query>
query Settings {
  metadata {
    settings {
      title
      description
      dropdownOptions {
        url
        latest
      }
    }
  }
}
</static-query>

<script>
import Logo from "@/components/Logo";
import {
  ArrowRightCircleIcon,
  ZapIcon,
  CodeIcon,
  AwardIcon,
  ClipboardIcon,
} from "vue-feather-icons";

export default {
  components: {
    Logo,
    ArrowRightCircleIcon,
    ZapIcon,
    CodeIcon,
    AwardIcon,
    ClipboardIcon,
  },

  metaInfo() {
    const title = this.$static.metadata.settings.title;
    const description = this.$static.metadata.settings.description;

    return {
      title: title,
      meta: [
        {
          name: "description",
          content: description,
        },
        {
          key: "og:title",
          name: "og:title",
          content: title,
        },
        {
          key: "twitter:title",
          name: "twitter:title",
          content: title,
        },
        {
          key: "og:description",
          name: "og:description",
          content: description,
        },
        {
          key: "twitter:description",
          name: "twitter:description",
          content: description,
        },
      ],
    };
  },

  methods: {
    latestDocs() {
      let options = this.$static.metadata.settings.dropdownOptions.find(
        (o) => o.latest === true
      );

      if (options) {
        return options.url;
      }
    },
  },
};
</script>

<style>
.home-links a {
  margin-right: 1rem;
}
</style>
