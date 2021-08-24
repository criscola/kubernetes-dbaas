"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[81],{3905:function(e,t,n){n.d(t,{Zo:function(){return c},kt:function(){return m}});var a=n(7294);function r(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function o(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){r(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function s(e,t){if(null==e)return{};var n,a,r=function(e,t){if(null==e)return{};var n,a,r={},i=Object.keys(e);for(a=0;a<i.length;a++)n=i[a],t.indexOf(n)>=0||(r[n]=e[n]);return r}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(a=0;a<i.length;a++)n=i[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(r[n]=e[n])}return r}var l=a.createContext({}),u=function(e){var t=a.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):o(o({},t),e)),n},c=function(e){var t=u(e.components);return a.createElement(l.Provider,{value:t},e.children)},d={inlineCode:"code",wrapper:function(e){var t=e.children;return a.createElement(a.Fragment,{},t)}},p=a.forwardRef((function(e,t){var n=e.components,r=e.mdxType,i=e.originalType,l=e.parentName,c=s(e,["components","mdxType","originalType","parentName"]),p=u(n),m=r,h=p["".concat(l,".").concat(m)]||p[m]||d[m]||i;return n?a.createElement(h,o(o({ref:t},c),{},{components:n})):a.createElement(h,o({ref:t},c))}));function m(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var i=n.length,o=new Array(i);o[0]=p;var s={};for(var l in t)hasOwnProperty.call(t,l)&&(s[l]=t[l]);s.originalType=e,s.mdxType="string"==typeof e?e:r,o[1]=s;for(var u=2;u<i;u++)o[u]=n[u];return a.createElement.apply(null,o)}return a.createElement.apply(null,n)}p.displayName="MDXCreateElement"},3919:function(e,t,n){function a(e){return!0===/^(\w*:|\/\/)/.test(e)}function r(e){return void 0!==e&&!a(e)}n.d(t,{b:function(){return a},Z:function(){return r}})},4996:function(e,t,n){n.d(t,{C:function(){return i},Z:function(){return o}});var a=n(2263),r=n(3919);function i(){var e=(0,a.Z)().siteConfig,t=(e=void 0===e?{}:e).baseUrl,n=void 0===t?"/":t,i=e.url;return{withBaseUrl:function(e,t){return function(e,t,n,a){var i=void 0===a?{}:a,o=i.forcePrependBaseUrl,s=void 0!==o&&o,l=i.absolute,u=void 0!==l&&l;if(!n)return n;if(n.startsWith("#"))return n;if((0,r.b)(n))return n;if(s)return t+n;var c=n.startsWith(t)?n:t+n.replace(/^\//,"");return u?e+c:c}(i,n,e,t)}}}function o(e,t){return void 0===t&&(t={}),(0,i().withBaseUrl)(e,t)}},9228:function(e,t,n){n.r(t),n.d(t,{frontMatter:function(){return l},contentTitle:function(){return u},metadata:function(){return c},toc:function(){return d},default:function(){return m}});var a=n(7462),r=n(3366),i=(n(7294),n(3905)),o=n(4996),s=["components"],l={sidebar_position:1},u="Overview",c={unversionedId:"overview",id:"overview",isDocsHomePage:!1,title:"Overview",description:"License",source:"@site/docs/overview.mdx",sourceDirName:".",slug:"/overview",permalink:"/kubernetes-dbaas/docs/overview",editUrl:"https://github.com/bedag/kubernetes-dbaas/edit/main/website/docs/overview.mdx",version:"current",sidebarPosition:1,frontMatter:{sidebar_position:1},sidebar:"tutorialSidebar",next:{title:"Prerequisites",permalink:"/kubernetes-dbaas/docs/dbms-configuration/prerequisites"}},d=[{value:"Description",id:"description",children:[]},{value:"Why?",id:"why",children:[{value:"In brief",id:"in-brief",children:[]},{value:"Kubernetes operators",id:"kubernetes-operators",children:[]},{value:"Goals",id:"goals",children:[]}]},{value:"Main features",id:"main-features",children:[]},{value:"Concepts",id:"concepts",children:[{value:"Custom resources",id:"custom-resources",children:[]},{value:"Operations",id:"operations",children:[]},{value:"An example",id:"an-example",children:[]}]},{value:"Supported DBMS",id:"supported-dbms",children:[]},{value:"Contributing",id:"contributing",children:[]}],p={toc:d};function m(e){var t=e.components,l=(0,r.Z)(e,s);return(0,i.kt)("wrapper",(0,a.Z)({},p,l,{components:t,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"overview"},"Overview"),(0,i.kt)("p",null,(0,i.kt)("a",{parentName:"p",href:"https://opensource.org/licenses/Apache-2.0"},(0,i.kt)("img",{parentName:"a",src:"https://img.shields.io/badge/License-Apache%202.0-blue.svg",alt:"License"})),"\n",(0,i.kt)("a",{parentName:"p",href:"https://pkg.go.dev/github.com/bedag/kubernetes-dbaas"},(0,i.kt)("img",{parentName:"a",src:"https://pkg.go.dev/badge/github.com/bedag/kubernetes-dbaas.svg",alt:"Go Reference"})),"\n",(0,i.kt)("a",{parentName:"p",href:"https://goreportcard.com/report/github.com/bedag/kubernetes-dbaas"},(0,i.kt)("img",{parentName:"a",src:"https://goreportcard.com/badge/github.com/bedag/kubernetes-dbaas",alt:"Go Report Card"})),"\n",(0,i.kt)("a",{parentName:"p",href:"https://artifacthub.io/packages/helm/kubernetes-dbaas/kubernetes-dbaas"},(0,i.kt)("img",{parentName:"a",src:"https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/kubernetes-dbaas",alt:"Artifact Hub"})),"\n",(0,i.kt)("a",{parentName:"p",href:"https://github.com/bedag/kubernetes-dbaas/actions/workflows/go.yaml"},(0,i.kt)("img",{parentName:"a",src:"https://github.com/bedag/kubernetes-dbaas/actions/workflows/go.yaml/badge.svg",alt:"Test Suite"}))),(0,i.kt)("div",{class:"markdown-centered"},(0,i.kt)("img",{alt:"Kubernetes DBaaS Logo",src:(0,o.Z)("/img/logo.svg"),width:"40%"})),(0,i.kt)("h2",{id:"description"},"Description"),(0,i.kt)("p",null,"The ",(0,i.kt)("strong",{parentName:"p"},"Kubernetes Database-as-a-Service (DBaaS) Operator"),' ("the Operator") is a ',(0,i.kt)("a",{parentName:"p",href:"https://kubernetes.io/docs/concepts/extend-kubernetes/operator/"},"Kubernetes\nOperator")," used\nto provision database instances in database management systems:"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},"The Operator can be easily configured and installed in a Kubernetes cluster\nusing the provided Helm Chart."),(0,i.kt)("li",{parentName:"ul"},"End-users such as software developers are able to create new database\ninstances by writing simple Database custom resources. "),(0,i.kt)("li",{parentName:"ul"},"Operations on DBMS are implemented using stored procedures called by the\nOperator whenever necessary, allowing you to define your own custom logic.   "),(0,i.kt)("li",{parentName:"ul"},"Credentials to access provisioned database instances are saved into Kubernetes\nSecrets.")),(0,i.kt)("p",null,"Written using the Go programming language."),(0,i.kt)("h2",{id:"why"},"Why?"),(0,i.kt)("h3",{id:"in-brief"},"In brief"),(0,i.kt)("p",null,"There are cases where an organization cannot or does not want to host their critical\ndata in cloud or distributed environments, and searches for a way to bridge the\ngap between their Kubernetes clusters and Database Management System (DBMS)\nsolutions. Medium to large organizations are often composed by distinct\nprofessional figures such as software developers, system administrators and\ndatabase administrators, each with its own need: "),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("strong",{parentName:"li"},"Software developers")," (end-users) would like to have their DB instances\nprovisioned as soon as possible using a user-friendly interface. "),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("strong",{parentName:"li"},"System administrators")," (sysadmins) would like to have a flexible, declarative\nsolution that is well-integrated in the Kubernetes ecosystem. "),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("strong",{parentName:"li"},"Database administrators")," (DBAs) must retain control over the core business\nlogic behind database provisioning as much as possible while still automating\nthe process in order to save time. ")),(0,i.kt)("h3",{id:"kubernetes-operators"},"Kubernetes operators"),(0,i.kt)("p",null,"This is where an on-premise Database-as-a-Service can help to satisfy those\nneeds. Kubernetes offers an innovative way to extend it by creating a\nKubernetes operator. An operator is a specific pattern used to handle the life\ncycle of Kubernetes resources. Its goal is to capture the human natural way of\nperforming tasks in order to automate processes that would otherwise be\ncarried out manually. Due to the infinite number of possibilities when\ndeploying and administering an application, Kubernetes can be extended using\nthe operator pattern with the intention of encapsulating complex business\nlogic, such as interacting with external services and performing a serie of\ntasks."),(0,i.kt)("h3",{id:"goals"},"Goals"),(0,i.kt)("p",null,"One of the Operator's strongest goals is having a clear ",(0,i.kt)("strong",{parentName:"p"},"separation of\nconcerns")," between end-users, sysadmins and DBAs: DBAs can retain full control\non the life cycle of database instances by creating stored procedures or an\nequivalent mechanism for each operation. This decouples the configuration from\nthe implementation and ensures a well-defined boundary between the Kubernetes\nand database worlds. ",(0,i.kt)("strong",{parentName:"p"},"Companies with strict compliance requirements can configure\nan opaque provisioning system for databases where data and business logic is\nkept as close as possible to their location")," without having to resort to a\nmanaged service; the only requirement is a formal specification under form of\nKubernetes resources between the system and database infrastructures which\nprovide the Operator with the minimal amount of information needed to\ncommunicate with each supported DBMS."),(0,i.kt)("div",{className:"admonition admonition-info alert alert--info"},(0,i.kt)("div",{parentName:"div",className:"admonition-heading"},(0,i.kt)("h5",{parentName:"div"},(0,i.kt)("span",{parentName:"h5",className:"admonition-icon"},(0,i.kt)("svg",{parentName:"span",xmlns:"http://www.w3.org/2000/svg",width:"14",height:"16",viewBox:"0 0 14 16"},(0,i.kt)("path",{parentName:"svg",fillRule:"evenodd",d:"M7 2.3c3.14 0 5.7 2.56 5.7 5.7s-2.56 5.7-5.7 5.7A5.71 5.71 0 0 1 1.3 8c0-3.14 2.56-5.7 5.7-5.7zM7 1C3.14 1 0 4.14 0 8s3.14 7 7 7 7-3.14 7-7-3.14-7-7-7zm1 3H6v5h2V4zm0 6H6v2h2v-2z"}))),"info")),(0,i.kt)("div",{parentName:"div",className:"admonition-content"},(0,i.kt)("p",{parentName:"div"},"The Operator can be used with database management systems hosted both inside and\noutside a Kubernetes cluster transparently."))),(0,i.kt)("h2",{id:"main-features"},"Main features"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},"Modern tech-stack, seamless Kubernetes integration"),(0,i.kt)("li",{parentName:"ul"},"Level-based logging, event recording, metrics, health/readiness probes..."),(0,i.kt)("li",{parentName:"ul"},"Flexible and powerful configuration"),(0,i.kt)("li",{parentName:"ul"},"Credential rotation"),(0,i.kt)("li",{parentName:"ul"},"Helm deployment"),(0,i.kt)("li",{parentName:"ul"},"Rate-limited requests")),(0,i.kt)("h2",{id:"concepts"},"Concepts"),(0,i.kt)("h3",{id:"custom-resources"},"Custom resources"),(0,i.kt)("p",null,"The Operator brings ",(0,i.kt)("strong",{parentName:"p"},"3 new custom resources")," into the cluster:"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("inlineCode",{parentName:"li"},"Database")," resources are used to describe Database instances.  "),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("inlineCode",{parentName:"li"},"DatabaseClass")," resources describe the format of the operations to be executed\non DB systems, what driver should be used to call operations and how data\nshould be saved into ",(0,i.kt)("inlineCode",{parentName:"li"},"Secrets"),"."),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("inlineCode",{parentName:"li"},"OperatorConfig")," is like a specialized ",(0,i.kt)("inlineCode",{parentName:"li"},"ConfigMap")," used to configure the\nOperator depending on the needs of the user organization. It contains also the\nlist of DBMS endpoints, including their DSN, bindings them to a particular\n",(0,i.kt)("inlineCode",{parentName:"li"},"DatabaseClass"),".")),(0,i.kt)("span",{class:"markdown-centered"},(0,i.kt)("img",{alt:"Custom resources",src:(0,o.Z)("/img/diagrams/custom-resources.svg")})),(0,i.kt)("h3",{id:"operations"},"Operations"),(0,i.kt)("p",null,"There are currently ",(0,i.kt)("strong",{parentName:"p"},"3 operations")," supported by the Operator:"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},"Database creation"),(0,i.kt)("li",{parentName:"ul"},"Database deletion"),(0,i.kt)("li",{parentName:"ul"},"Database credential rotation")),(0,i.kt)("p",null,"The ",(0,i.kt)("a",{parentName:"p",href:"https://kubernetes.io/docs/concepts/architecture/controller"},"control loop")," of the Operator can be summarized by\nmeans of the following flowchart:"),(0,i.kt)("p",null,(0,i.kt)("img",{alt:"System diagram",src:n(3831).Z})),(0,i.kt)("h3",{id:"an-example"},"An example"),(0,i.kt)("p",null,"The following diagram shows what happens when an operation is executed on\na ",(0,i.kt)("inlineCode",{parentName:"p"},"Database")," resource:"),(0,i.kt)("p",null,(0,i.kt)("img",{alt:"System diagram",src:n(8879).Z})),(0,i.kt)("ol",null,(0,i.kt)("li",{parentName:"ol"},"The Operator watches the cluster for a new event generated by a Database\nresource, i.e. creation, deletion or credential rotation."),(0,i.kt)("li",{parentName:"ol"},"The Operator calls the relative stored procedure or equivalent mechanism on\nthe DBMS."),(0,i.kt)("li",{parentName:"ol"},"The DBMS executes the stored procedure according to the implementation of\nthe database administrator."),(0,i.kt)("li",{parentName:"ol"},"Finally, the Operator acts on the Secret by creating, deleting or updating it\nwith the data returned by the operation.")),(0,i.kt)("h2",{id:"supported-dbms"},"Supported DBMS"),(0,i.kt)("p",null,"See ",(0,i.kt)("a",{parentName:"p",href:"/docs/dbms-configuration/prerequisites#supported-dbms"},"Supported DBMS"),"."),(0,i.kt)("h2",{id:"contributing"},"Contributing"),(0,i.kt)("p",null,"See ",(0,i.kt)("a",{parentName:"p",href:"/docs/contributing/how-to-contribute"},"Contributing introduction"),"."))}m.isMDXComponent=!0},8879:function(e,t,n){t.Z=n.p+"assets/images/01_system_diagram-a731b2a4ed6b466b39f84d722cfe59a1.png"},3831:function(e,t,n){t.Z=n.p+"assets/images/01_system_flowchart_diagram-fecc16918b61fee640fa7c81a5bf930f.png"}}]);