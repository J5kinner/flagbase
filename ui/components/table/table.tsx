import React from 'react';
import { Table as AntdTable, TableColumnProps, TableProps as AntdTableProps, TableColumnType } from 'antd';
import styled from '@emotion/styled';

export type TableProps = {
  loading: boolean;
  dataSource: AntdTableProps<string>[];
  columns: TableColumnProps<TableColumnType<string>>[];
};

const StyledTable = styled(AntdTable)`
 .ant-table-thead > tr > th {
     background-color: white;
 }
`;

const Table: React.FC<TableProps> = ({ loading, columns, dataSource }) => {
  return <StyledTable loading={loading} dataSource={dataSource} columns={columns}>
  </StyledTable>;
};

export default Table;
