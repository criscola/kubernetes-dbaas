import React from 'react';
import clsx from 'clsx';
import styles from './HomepageFeatures.module.css';

const FeatureList = [
  {
    title: 'Use familiar tools',
    Svg: require('../../static/img/undraw_use_familiar_tools.svg').default,
    description: (
      <>
        DBAs implement stored procedures, cluster administrators deploy the Operator, end-users write short declarative YAML files.
      </>
    ),
  },
  {
    title: 'Simple self-service mechanism',
    Svg: require('../../static/img/undraw_separation_concerns2.svg').default,
    description: (
      <>
        You don't need to be a Kubernetes expert; once set up it's as simple as running &nbsp;
        <code>kubectl apply -f db.yaml</code>.
      </>
    ),
  },
  {
    title: 'Flexible configuration',
    Svg: require('../../static/img/undraw_flexibility.svg').default,
    description: (
      <>
        Designed with flexibility in mind from the ground-up, <strong>you</strong> are in charge of how databases are provisioned.
      </>
    ),
  },
    {
        title: 'Modern tech-stack',
        Svg: require('../../static/img/undraw_modern.svg').default,
        description: (
            <>
              Extend the Kubernetes API with custom resources and deploy the Operator using Helm.
            </>
        ),
    },
    {
        title: 'Open-source software',
        Svg: require('../../static/img/undraw_open_source.svg').default,
        description: (
            <>
              Fully deploy the Operator in your own infrastructure. It's open-source, and it will always be.
            </>
        ),
    },
    {
        title: 'Kubernetes Secrets',
        Svg: require('../../static/img/undraw_secrets.svg').default,
        description: (
            <>
              Database secrets
              are stored using Kubernetes Secrets, giving applications transparent, secure access to credentials.
            </>
        ),
    },
];

function Feature({Svg, title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <Svg className={styles.featureSvg} alt={title} />
      </div>
      <div className="text--center padding-horiz--md">
        <h3>{title}</h3>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
