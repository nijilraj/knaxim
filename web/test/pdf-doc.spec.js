import { shallowMount } from '@vue/test-utils'
import PdfDoc from '@/components/pdf/pdf-doc'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(PdfDoc, {
    stubs: ['b-col', 'b-row', 'b-container'],
    propsData: {
      ...options.props
    },
    methods: {
      ...options.methods
    },
    computed: {
      url () {
        return '<html></html>'
      },
      ...options.computed
    }
  })
}

describe('PdfDoc', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(PdfDoc)).toBe(true)
  })
})