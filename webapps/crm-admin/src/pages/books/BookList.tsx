import PageCard from "@/components/PageCard";
import { Table } from "react-bootstrap";

function BookList() {
  return (
    <PageCard>
      <PageCard.Header>
        {/* 查询条件：图书名、作者、分类、出版社、出版时间 */}
        {/* 按钮功能：查询、重置、新增、导入、批量删除 */}
        {/* 显示字段： */}
        {/* 排序支持：添加时间、出版时间、价格、书名、作者、 */}
        搜索内容
      </PageCard.Header>
      <PageCard.Body>
        <Table striped bordered hover>
          <thead>
            <tr>
              <th>#</th>
              <th>First Name</th>
              <th>Last Name</th>
              <th>Username</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>1</td>
              <td>Mark</td>
              <td>Otto</td>
              <td>@mdo</td>
            </tr>
            <tr>
              <td>2</td>
              <td>Jacob</td>
              <td>Thornton</td>
              <td>@fat</td>
            </tr>
            <tr>
              <td>3</td>
              <td colSpan={2}>Larry the Bird</td>
              <td>@twitter</td>
            </tr>
          </tbody>
        </Table>
      </PageCard.Body>
    </PageCard>
  );
}

export default BookList;
