export default {
  "title": "Kubernetes Database-as-a-Service",
  "tagline": "A unique Kubernetes Database-as-a-Service (DBaaS) Operator for declarative, self-service database provisioning in DBMS solutions.",
  "url": "https://criscola.github.io/kubernetes-dbaas",
  "baseUrl": "/kubernetes-dbaas/",
  "favicon": "/img/favicon.ico",
  "trailingSlash": false,
  "onBrokenLinks": "throw",
  "onBrokenMarkdownLinks": "warn",
  "organizationName": "criscola",
  "projectName": "kubernetes-dbaas",
  "themeConfig": {
    "navbar": {
      "title": "Home",
      "logo": {
        "alt": "Kubernetes Database-as-a-Service logo",
        "src": "img/logo.svg"
      },
      "items": [
        {
          "type": "doc",
          "docId": "overview",
          "position": "left",
          "label": "Docs",
          "activeSidebarClassName": "navbar__link--active"
        },
        {
          "href": "https://github.com/bedag/kubernetes-dbaas",
          "label": "GitHub",
          "position": "right"
        }
      ],
      "hideOnScroll": false
    },
    "footer": {
      "style": "dark",
      "links": [
        {
          "title": "Docs",
          "items": [
            {
              "label": "Documentation",
              "to": "/docs/overview"
            },
            {
              "label": "Godocs",
              "to": "https://pkg.go.dev/github.com/bedag/kubernetes-dbaas"
            },
            {
              "label": "Contributing",
              "to": "/docs/contributing/how-to-contribute"
            },
            {
              "label": "Legals",
              "to": "/docs/legals"
            }
          ]
        },
        {
          "title": "Links",
          "items": [
            {
              "label": "Source code",
              "href": "https://github.com/bedag/kubernetes-dbaas"
            },
            {
              "label": "Issue tracker",
              "href": "https://github.com/bedag/kubernetes-dbaas/issues"
            },
            {
              "label": "DockerHub repository",
              "href": "https://hub.docker.com/r/bedag/kubernetes-dbaas"
            },
            {
              "label": "ArtifactHub repository",
              "href": "#"
            }
          ]
        },
        {
          "title": "More",
          "items": [
            {
              "label": "Bedag Informatik AG",
              "href": "https://www.bedag.ch"
            },
            {
              "label": "LinkedIn",
              "href": "https://www.linkedin.com/company/bedag"
            },
            {
              "label": "Instagram",
              "href": "https://www.instagram.com/bedaginformatik"
            }
          ]
        }
      ],
      "copyright": "Copyright © 2021 <a href='https://github.com/criscola'>Cristiano Colangelo</a> — Developed for Bedag Informatik AG. <br>"
    },
    "prism": {
      "theme": {
        "plain": {
          "color": "#393A34",
          "backgroundColor": "#f6f8fa"
        },
        "styles": [
          {
            "types": [
              "comment",
              "prolog",
              "doctype",
              "cdata"
            ],
            "style": {
              "color": "#999988",
              "fontStyle": "italic"
            }
          },
          {
            "types": [
              "namespace"
            ],
            "style": {
              "opacity": 0.7
            }
          },
          {
            "types": [
              "string",
              "attr-value"
            ],
            "style": {
              "color": "#e3116c"
            }
          },
          {
            "types": [
              "punctuation",
              "operator"
            ],
            "style": {
              "color": "#393A34"
            }
          },
          {
            "types": [
              "entity",
              "url",
              "symbol",
              "number",
              "boolean",
              "variable",
              "constant",
              "property",
              "regex",
              "inserted"
            ],
            "style": {
              "color": "#36acaa"
            }
          },
          {
            "types": [
              "atrule",
              "keyword",
              "attr-name",
              "selector"
            ],
            "style": {
              "color": "#00a4db"
            }
          },
          {
            "types": [
              "function",
              "deleted",
              "tag"
            ],
            "style": {
              "color": "#d73a49"
            }
          },
          {
            "types": [
              "function-variable"
            ],
            "style": {
              "color": "#6f42c1"
            }
          },
          {
            "types": [
              "tag",
              "selector",
              "keyword"
            ],
            "style": {
              "color": "#00009f"
            }
          }
        ]
      },
      "darkTheme": {
        "plain": {
          "color": "#F8F8F2",
          "backgroundColor": "#282A36"
        },
        "styles": [
          {
            "types": [
              "prolog",
              "constant",
              "builtin"
            ],
            "style": {
              "color": "rgb(189, 147, 249)"
            }
          },
          {
            "types": [
              "inserted",
              "function"
            ],
            "style": {
              "color": "rgb(80, 250, 123)"
            }
          },
          {
            "types": [
              "deleted"
            ],
            "style": {
              "color": "rgb(255, 85, 85)"
            }
          },
          {
            "types": [
              "changed"
            ],
            "style": {
              "color": "rgb(255, 184, 108)"
            }
          },
          {
            "types": [
              "punctuation",
              "symbol"
            ],
            "style": {
              "color": "rgb(248, 248, 242)"
            }
          },
          {
            "types": [
              "string",
              "char",
              "tag",
              "selector"
            ],
            "style": {
              "color": "rgb(255, 121, 198)"
            }
          },
          {
            "types": [
              "keyword",
              "variable"
            ],
            "style": {
              "color": "rgb(189, 147, 249)",
              "fontStyle": "italic"
            }
          },
          {
            "types": [
              "comment"
            ],
            "style": {
              "color": "rgb(98, 114, 164)"
            }
          },
          {
            "types": [
              "attr-name"
            ],
            "style": {
              "color": "rgb(241, 250, 140)"
            }
          }
        ]
      },
      "additionalLanguages": []
    },
    "colorMode": {
      "defaultMode": "light",
      "disableSwitch": false,
      "respectPrefersColorScheme": false,
      "switchConfig": {
        "darkIcon": "🌜",
        "darkIconStyle": {},
        "lightIcon": "🌞",
        "lightIconStyle": {}
      }
    },
    "docs": {
      "versionPersistence": "localStorage"
    },
    "metadatas": [],
    "hideableSidebar": false
  },
  "presets": [
    [
      "@docusaurus/preset-classic",
      {
        "docs": {
          "sidebarPath": "/home/runner/work/kubernetes-dbaas/kubernetes-dbaas/website/sidebars.js",
          "editUrl": "https://github.com/bedag/kubernetes-dbaas/edit/main/website"
        },
        "theme": {
          "customCss": "/home/runner/work/kubernetes-dbaas/kubernetes-dbaas/website/src/css/custom.css"
        }
      }
    ]
  ],
  "baseUrlIssueBanner": true,
  "i18n": {
    "defaultLocale": "en",
    "locales": [
      "en"
    ],
    "localeConfigs": {}
  },
  "onDuplicateRoutes": "warn",
  "customFields": {},
  "plugins": [],
  "themes": [],
  "titleDelimiter": "|",
  "noIndex": false
};