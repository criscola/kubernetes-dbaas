
import React from 'react';
import ComponentCreator from '@docusaurus/ComponentCreator';

export default [
  {
    path: '/kubernetes-dbaas/',
    component: ComponentCreator('/kubernetes-dbaas/','9a1'),
    exact: true
  },
  {
    path: '/kubernetes-dbaas/docs',
    component: ComponentCreator('/kubernetes-dbaas/docs','6d3'),
    routes: [
      {
        path: '/kubernetes-dbaas/docs/contributing/architecture',
        component: ComponentCreator('/kubernetes-dbaas/docs/contributing/architecture','a2a'),
        exact: true,
        'sidebar': "tutorialSidebar"
      },
      {
        path: '/kubernetes-dbaas/docs/contributing/branching-standard',
        component: ComponentCreator('/kubernetes-dbaas/docs/contributing/branching-standard','6df'),
        exact: true,
        'sidebar': "tutorialSidebar"
      },
      {
        path: '/kubernetes-dbaas/docs/contributing/ci',
        component: ComponentCreator('/kubernetes-dbaas/docs/contributing/ci','580'),
        exact: true,
        'sidebar': "tutorialSidebar"
      },
      {
        path: '/kubernetes-dbaas/docs/contributing/how-to-contribute',
        component: ComponentCreator('/kubernetes-dbaas/docs/contributing/how-to-contribute','c9b'),
        exact: true,
        'sidebar': "tutorialSidebar"
      },
      {
        path: '/kubernetes-dbaas/docs/contributing/testing',
        component: ComponentCreator('/kubernetes-dbaas/docs/contributing/testing','af7'),
        exact: true,
        'sidebar': "tutorialSidebar"
      },
      {
        path: '/kubernetes-dbaas/docs/dbms-configuration/operations',
        component: ComponentCreator('/kubernetes-dbaas/docs/dbms-configuration/operations','50a'),
        exact: true,
        'sidebar': "tutorialSidebar"
      },
      {
        path: '/kubernetes-dbaas/docs/dbms-configuration/prerequisites',
        component: ComponentCreator('/kubernetes-dbaas/docs/dbms-configuration/prerequisites','7ca'),
        exact: true,
        'sidebar': "tutorialSidebar"
      },
      {
        path: '/kubernetes-dbaas/docs/dbms-configuration/samples',
        component: ComponentCreator('/kubernetes-dbaas/docs/dbms-configuration/samples','972'),
        exact: true,
        'sidebar': "tutorialSidebar"
      },
      {
        path: '/kubernetes-dbaas/docs/legals',
        component: ComponentCreator('/kubernetes-dbaas/docs/legals','07b'),
        exact: true,
        'sidebar': "tutorialSidebar"
      },
      {
        path: '/kubernetes-dbaas/docs/operator-configuration/cli-arguments',
        component: ComponentCreator('/kubernetes-dbaas/docs/operator-configuration/cli-arguments','de2'),
        exact: true,
        'sidebar': "tutorialSidebar"
      },
      {
        path: '/kubernetes-dbaas/docs/operator-configuration/credential-rotation',
        component: ComponentCreator('/kubernetes-dbaas/docs/operator-configuration/credential-rotation','4ad'),
        exact: true,
        'sidebar': "tutorialSidebar"
      },
      {
        path: '/kubernetes-dbaas/docs/operator-configuration/databaseclasses',
        component: ComponentCreator('/kubernetes-dbaas/docs/operator-configuration/databaseclasses','e89'),
        exact: true,
        'sidebar': "tutorialSidebar"
      },
      {
        path: '/kubernetes-dbaas/docs/operator-configuration/logging-monitoring',
        component: ComponentCreator('/kubernetes-dbaas/docs/operator-configuration/logging-monitoring','687'),
        exact: true,
        'sidebar': "tutorialSidebar"
      },
      {
        path: '/kubernetes-dbaas/docs/operator-configuration/main-configuration',
        component: ComponentCreator('/kubernetes-dbaas/docs/operator-configuration/main-configuration','e90'),
        exact: true,
        'sidebar': "tutorialSidebar"
      },
      {
        path: '/kubernetes-dbaas/docs/operator-configuration/prerequisites',
        component: ComponentCreator('/kubernetes-dbaas/docs/operator-configuration/prerequisites','02a'),
        exact: true,
        'sidebar': "tutorialSidebar"
      },
      {
        path: '/kubernetes-dbaas/docs/operator-configuration/tips-and-tricks',
        component: ComponentCreator('/kubernetes-dbaas/docs/operator-configuration/tips-and-tricks','b9e'),
        exact: true,
        'sidebar': "tutorialSidebar"
      },
      {
        path: '/kubernetes-dbaas/docs/operator-deployment/development-deployment',
        component: ComponentCreator('/kubernetes-dbaas/docs/operator-deployment/development-deployment','ef5'),
        exact: true,
        'sidebar': "tutorialSidebar"
      },
      {
        path: '/kubernetes-dbaas/docs/operator-deployment/helm',
        component: ComponentCreator('/kubernetes-dbaas/docs/operator-deployment/helm','ba9'),
        exact: true,
        'sidebar': "tutorialSidebar"
      },
      {
        path: '/kubernetes-dbaas/docs/operator-deployment/vanilla-deployment',
        component: ComponentCreator('/kubernetes-dbaas/docs/operator-deployment/vanilla-deployment','c77'),
        exact: true,
        'sidebar': "tutorialSidebar"
      },
      {
        path: '/kubernetes-dbaas/docs/overview',
        component: ComponentCreator('/kubernetes-dbaas/docs/overview','c91'),
        exact: true,
        'sidebar': "tutorialSidebar"
      },
      {
        path: '/kubernetes-dbaas/docs/usage',
        component: ComponentCreator('/kubernetes-dbaas/docs/usage','e4e'),
        exact: true,
        'sidebar': "tutorialSidebar"
      }
    ]
  },
  {
    path: '*',
    component: ComponentCreator('*')
  }
];
