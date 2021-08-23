const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

/** @type {import('@docusaurus/types').DocusaurusConfig} */
module.exports = {
  title: 'Kubernetes Database-as-a-Service',
  tagline: 'A unique Kubernetes Database-as-a-Service (DBaaS) Operator for declarative, self-service database' +
    ' provisioning in DBMS' +
    ' solutions.',
  url: 'https://criscola.github.io',
  baseUrl: '/kubernetes-dbaas/',
  favicon: '/img/favicon.ico',
  trailingSlash: false,
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  organizationName: 'criscola', // Usually your GitHub org/user name.
  projectName: 'kubernetes-dbaas', // Usually your repo name.
  themeConfig: {
    navbar: {
      title: 'Home',
      logo: {
        alt: 'Kubernetes Database-as-a-Service logo',
        src: 'img/logo.svg',
      },
      items: [
        {
          type: 'doc',
          docId: 'overview',
          position: 'left',
          label: 'Docs',
        },
        {
          href: 'https://github.com/bedag/kubernetes-dbaas',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            {
              label: 'Documentation',
              to: '/docs/overview',
            },            
            {
              label: 'Godocs',
              to: 'https://pkg.go.dev/github.com/bedag/kubernetes-dbaas',
            },

            {
              label: 'Contributing',
              to: '/docs/contributing/how-to-contribute',
            },
            {
              label: 'Legals',
              to: '/docs/legals',
            },
          ],
        },
        {
          title: 'Links',
          items: [
            {
              label: 'Source code',
              href: 'https://github.com/bedag/kubernetes-dbaas',
            },
            {
              label: 'Issue tracker',
              href: 'https://github.com/bedag/kubernetes-dbaas/issues',
            },
            {
              label: 'DockerHub repository',
              href: 'https://hub.docker.com/r/bedag/kubernetes-dbaas',
            },
            {
              label: 'ArtifactHub repository',
              href: '#',
            },
          ],
        },
        {
          title: 'More',
          items: [
            {
              label: 'Bedag Informatik AG',
              href: 'https://www.bedag.ch',
            },
            {
              label: 'LinkedIn',
              href: 'https://www.linkedin.com/company/bedag',
            },
            {
              label: 'Instagram',
              href: 'https://www.instagram.com/bedaginformatik',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} <a href='https://github.com/criscola'>Cristiano Colangelo</a> — Developed for Bedag Informatik AG. <br>`,
    },
    prism: {
      theme: lightCodeTheme,
      darkTheme: darkCodeTheme,
    },
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          editUrl:
            'https://github.com/bedag/kubernetes-dbaas/edit/main/website',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      },
    ],
  ]
};
