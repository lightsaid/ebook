import React from "react";
import "./style.scss";

const PageCard = ({ children }: { children: React.ReactNode }) => {
  return <div>{children}</div>;
};

PageCard.Header = ({ children }: { children: React.ReactNode }) =>  {
  return <div className="page-card-head">
    {children}
  </div>;
};

PageCard.Body = ({ children }: { children: React.ReactNode }) =>  {
  return <div className="page-card-body">
    {children}
  </div>;
};


export default PageCard;
